package main

import (
	"crypto/subtle"
	"encoding/json"
	"html/template"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

type App struct {
	templates *template.Template
	cfg       Config
	ejabberd  *EjabberdClient
}

type Config struct {
	ListenAddr         string
	BasePath           string
	Domain             string
	AdminUser          string
	AdminPassword      string
	EjabberdAPI        string
	EjabberdAPIUser    string
	EjabberdAPIPass    string
	EjabberdConfigPath string
	Server             EjabberdConfigSnapshot
}

func main() {
	cfg := loadConfig()
	validateRuntimeConfig(cfg)

	tpl, err := template.New("root").Funcs(template.FuncMap{
		"tr":   tr,
		"join": strings.Join,
		"p":       func(path string) string { return joinBasePath(cfg.BasePath, path) },
		"safeURL": safeExternalURL,
	}).Parse(adminTemplate)
	if err != nil {
		log.Fatal(err)
	}

	httpClient := &http.Client{
		Timeout: 8 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	app := &App{
		templates: tpl,
		cfg:       cfg,
		ejabberd: &EjabberdClient{
			baseURL:  cfg.EjabberdAPI,
			username: cfg.EjabberdAPIUser,
			password: cfg.EjabberdAPIPass,
			client:   httpClient,
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/assets/material3.css", app.material3CSS)
	mux.HandleFunc("/assets/theme.js", app.themeJS)
	mux.HandleFunc("/healthz", app.health)
	mux.HandleFunc("/readyz", app.ready)
	mux.HandleFunc("/lang", app.setLanguage)
	mux.HandleFunc("/admin", app.adminHome)
	mux.HandleFunc("/admin/invites/create", app.adminCreateInvite)
	mux.HandleFunc("/admin/invites/revoke", app.adminRevokeInvite)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, joinBasePath(cfg.BasePath, "/admin"), http.StatusSeeOther)
	})

	var handler http.Handler = mux
	if cfg.BasePath != "" {
		stripped := http.StripPrefix(cfg.BasePath, mux)
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == cfg.BasePath {
				http.Redirect(w, r, cfg.BasePath+"/admin", http.StatusSeeOther)
				return
			}
			if !strings.HasPrefix(r.URL.Path, cfg.BasePath+"/") {
				http.NotFound(w, r)
				return
			}
			stripped.ServeHTTP(w, r)
		})
	}

	server := &http.Server{
		Addr:              cfg.ListenAddr,
		Handler:           securityHeaders(handler),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Printf("XMPP Admin listening on %s; host=%s; ejabberd_config=%s; api=%s",
		cfg.ListenAddr, cfg.Domain, cfg.EjabberdConfigPath, cfg.EjabberdAPI)
	log.Fatal(server.ListenAndServe())
}

func loadConfig() Config {
	configPath := resolveEjabberdConfigPath(os.Getenv("EJABBERD_CONFIG"))
	domainOverride := strings.TrimSpace(os.Getenv("XMPP_DOMAIN"))
	snapshot, configErr := loadEjabberdConfigForHost(configPath, domainOverride)

	domain := domainOverride
	if domain == "" && configErr == nil {
		domain = snapshot.PrimaryHost
	}
	if domain == "" {
		log.Fatalf("cannot determine XMPP domain: set XMPP_DOMAIN or provide readable ejabberd config at %s", configPath)
	}

	apiURL := strings.TrimRight(strings.TrimSpace(os.Getenv("EJABBERD_API")), "/")
	if apiURL == "" && configErr == nil && snapshot.HTTPAPIListenerFound {
		apiURL = strings.TrimRight(snapshot.HTTPAPIURL, "/")
	}
	if apiURL == "" {
		apiURL = "http://127.0.0.1:5281/api"
	}

	if configErr != nil {
		log.Printf("warning: ejabberd config was not loaded from %s: %v; using explicit/fallback settings", configPath, configErr)
		snapshot = EjabberdConfigSnapshot{Path: configPath}
	}

	return Config{
		ListenAddr:         env("ADMIN_LISTEN", "127.0.0.1:8090"),
		BasePath:           normalizeBasePath(env("ADMIN_BASE_PATH", "")),
		Domain:             domain,
		AdminUser:          env("ADMIN_USER", "admin"),
		AdminPassword:      os.Getenv("ADMIN_PASSWORD"),
		EjabberdAPI:        apiURL,
		EjabberdAPIUser:    os.Getenv("EJABBERD_API_USER"),
		EjabberdAPIPass:    os.Getenv("EJABBERD_API_PASSWORD"),
		EjabberdConfigPath: configPath,
		Server:             snapshot,
	}
}

func env(k, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(k)); v != "" {
		return v
	}
	return fallback
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=(), payment=(), usb=()")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self'; script-src 'self'; img-src 'self' data:; connect-src 'self'; form-action 'self'; base-uri 'none'; object-src 'none'; frame-ancestors 'none'")
		next.ServeHTTP(w, r)
	})
}

func (a *App) requireAdmin(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	u, p, ok := r.BasicAuth()
	userOK := subtle.ConstantTimeCompare([]byte(u), []byte(a.cfg.AdminUser)) == 1
	passOK := subtle.ConstantTimeCompare([]byte(p), []byte(a.cfg.AdminPassword)) == 1
	if !ok || !userOK || !passOK {
		w.Header().Set("WWW-Authenticate", `Basic realm="XMPP Admin", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}
	if r.Method == http.MethodPost && !sameOrigin(r) {
		http.Error(w, "Bad origin", http.StatusForbidden)
		return false
	}
	return true
}

func sameOrigin(r *http.Request) bool {
	if site := strings.ToLower(strings.TrimSpace(r.Header.Get("Sec-Fetch-Site"))); site != "" && site != "same-origin" && site != "none" {
		return false
	}
	expectedScheme := strings.ToLower(strings.TrimSpace(r.Header.Get("X-Forwarded-Proto")))
	if expectedScheme == "" {
		if r.TLS != nil {
			expectedScheme = "https"
		} else {
			expectedScheme = "http"
		}
	}
	checkURL := func(raw string) bool {
		u, err := url.Parse(raw)
		if err != nil || u.Host == "" {
			return false
		}
		return strings.EqualFold(u.Host, r.Host) && strings.EqualFold(u.Scheme, expectedScheme)
	}
	if origin := strings.TrimSpace(r.Header.Get("Origin")); origin != "" {
		return checkURL(origin)
	}
	if referer := strings.TrimSpace(r.Header.Get("Referer")); referer != "" {
		return checkURL(referer)
	}
	return strings.EqualFold(strings.TrimSpace(r.Header.Get("Sec-Fetch-Site")), "same-origin")
}

func (a *App) currentServerConfig() EjabberdConfigSnapshot {
	snapshot, err := loadEjabberdConfigForHost(a.cfg.EjabberdConfigPath, a.cfg.Domain)
	if err != nil {
		return a.cfg.Server
	}
	return snapshot
}

func (a *App) health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func (a *App) ready(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Forwarded-For") != "" || r.Header.Get("X-Real-IP") != "" {
		http.NotFound(w, r)
		return
	}
	ok := false
	serverCfg := a.currentServerConfig()
	if serverCfg.InvitesEnabled {
		if _, err := a.ejabberd.ListInvites(r.Context(), a.cfg.Domain); err == nil {
			ok = true
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	if !ok {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	_ = json.NewEncoder(w).Encode(map[string]bool{"ok": ok})
}

func (a *App) setLanguage(w http.ResponseWriter, r *http.Request) {
	lang := r.URL.Query().Get("lang")
	if lang != langRU && lang != langEN {
		lang = langEN
	}
	setLanguageCookie(w, r, lang)
	http.Redirect(w, r, joinBasePath(a.cfg.BasePath, safeNext(r.URL.Query().Get("next"))), http.StatusSeeOther)
}

type AdminPageData struct {
	Invites    []NativeInvite
	Domain     string
	Created    bool
	Lang       string
	APIError   string
	Server     EjabberdConfigSnapshot
	ConfigPath string
}

func (a *App) adminHome(w http.ResponseWriter, r *http.Request) {
	if !a.requireAdmin(w, r) {
		return
	}
	if r.URL.Path != "/admin" {
		http.NotFound(w, r)
		return
	}
	a.renderAdminPage(w, r)
}

func (a *App) renderAdminPage(w http.ResponseWriter, r *http.Request) {
	lang := requestLanguage(r)
	serverCfg := a.currentServerConfig()
	data := AdminPageData{
		Domain:     a.cfg.Domain,
		Created:    r.URL.Query().Get("created") == "1",
		Lang:       lang,
		Server:     serverCfg,
		ConfigPath: a.cfg.EjabberdConfigPath,
	}
	if !serverCfg.InvitesEnabled {
		data.APIError = tr(lang, "mod_invites_disabled")
	} else {
		invites, err := a.ejabberd.ListInvites(r.Context(), a.cfg.Domain)
		if err != nil {
			data.APIError = err.Error()
		} else {
			data.Invites = invites
		}
	}
	if err := a.templates.ExecuteTemplate(w, "admin", data); err != nil {
		log.Print(err)
	}
}

func (a *App) adminCreateInvite(w http.ResponseWriter, r *http.Request) {
	if !a.requireAdmin(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}
	lang := requestLanguage(r)
	if !a.currentServerConfig().InvitesEnabled {
		http.Error(w, tr(lang, "mod_invites_disabled"), http.StatusServiceUnavailable)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	username := normalizeUsername(r.FormValue("username"))
	if username != "" && !validUsername(username) {
		http.Error(w, tr(lang, "invalid_username"), http.StatusBadRequest)
		return
	}
	if _, err := a.ejabberd.GenerateInvite(r.Context(), a.cfg.Domain, username); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	http.Redirect(w, r, joinBasePath(a.cfg.BasePath, "/admin")+"?created=1", http.StatusSeeOther)
}

func (a *App) adminRevokeInvite(w http.ResponseWriter, r *http.Request) {
	if !a.requireAdmin(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	token := strings.TrimSpace(r.FormValue("token"))
	if token == "" || len(token) > 256 {
		http.Error(w, "bad token", http.StatusBadRequest)
		return
	}
	if err := a.ejabberd.ExpireInvite(r.Context(), a.cfg.Domain, token); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	http.Redirect(w, r, joinBasePath(a.cfg.BasePath, "/admin"), http.StatusSeeOther)
}

func normalizeUsername(s string) string {
	return strings.TrimSpace(s)
}

func validUsername(s string) bool {
	if s == "" || len(s) > 255 || !utf8.ValidString(s) {
		return false
	}
	if strings.ContainsAny(s, "@/") {
		return false
	}
	for _, r := range s {
		if unicode.IsControl(r) || unicode.IsSpace(r) {
			return false
		}
	}
	// ejabberd remains authoritative and applies its own nodeprep/JID validation.
	return true
}


func normalizeBasePath(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "/" {
		return ""
	}
	if !strings.HasPrefix(raw, "/") {
		raw = "/" + raw
	}
	return strings.TrimRight(raw, "/")
}

func joinBasePath(base, path string) string {
	if path == "" {
		path = "/"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if base == "" {
		return path
	}
	return base + path
}


func validateRuntimeConfig(cfg Config) {
	if cfg.AdminUser == "" || strings.Contains(cfg.AdminUser, ":") || containsCredentialControl(cfg.AdminUser) {
		log.Fatal("ADMIN_USER must be non-empty and must not contain ':' or control characters")
	}
	if len(cfg.AdminPassword) < 12 || containsCredentialControl(cfg.AdminPassword) {
		log.Fatal("ADMIN_PASSWORD must contain at least 12 characters and no CR/LF/NUL")
	}
	if (cfg.EjabberdAPIUser == "") != (cfg.EjabberdAPIPass == "") {
		log.Fatal("EJABBERD_API_USER and EJABBERD_API_PASSWORD must be configured together")
	}
	if cfg.EjabberdAPIUser != "" && (strings.Contains(cfg.EjabberdAPIUser, ":") || containsCredentialControl(cfg.EjabberdAPIUser) || containsCredentialControl(cfg.EjabberdAPIPass)) {
		log.Fatal("ejabberd API credentials contain unsupported control characters or ':' in the username")
	}
	u, err := url.Parse(cfg.EjabberdAPI)
	if err != nil || u.Hostname() == "" {
		log.Fatalf("invalid EJABBERD_API: %q", cfg.EjabberdAPI)
	}
	if u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		log.Fatal("EJABBERD_API must not contain URL credentials, query parameters, or fragments")
	}
	switch strings.ToLower(u.Scheme) {
	case "https":
		return
	case "http":
		host := strings.ToLower(u.Hostname())
		ip := net.ParseIP(host)
		if host == "localhost" || (ip != nil && ip.IsLoopback()) {
			return
		}
		log.Fatal("refusing plaintext EJABBERD_API over a non-loopback address; use HTTPS or a loopback API listener")
	default:
		log.Fatalf("unsupported EJABBERD_API scheme %q", u.Scheme)
	}
}


func safeExternalURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	if u.User != nil {
		return ""
	}
	switch strings.ToLower(u.Scheme) {
	case "http", "https":
		if u.Host == "" {
			return ""
		}
	case "xmpp":
		if !strings.HasPrefix(strings.ToLower(raw), "xmpp:") {
			return ""
		}
	default:
		return ""
	}
	return raw
}


func containsCredentialControl(value string) bool {
	for _, r := range value {
		if r == 0 || r == 10 || r == 13 {
			return true
		}
	}
	return false
}
