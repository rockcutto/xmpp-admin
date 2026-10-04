package main

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestLocaleCatalogsHaveSameKeys(t *testing.T) {
	en := messages[langEN]
	ru := messages[langRU]
	for key := range en {
		if _, ok := ru[key]; !ok {
			t.Errorf("missing RU translation for %q", key)
		}
	}
	for key := range ru {
		if _, ok := en[key]; !ok {
			t.Errorf("missing EN translation for %q", key)
		}
	}
}

func TestInvitePresentationHelpers(t *testing.T) {
	if got := formatTimestamp("2026-06-19T09:42:16Z"); got != "2026-06-19 09:42 UTC" {
		t.Fatalf("unexpected formatted timestamp: %q", got)
	}
	if got := inviteTypeLabel(langEN, "account_subscription"); got != "Account + contact" {
		t.Fatalf("unexpected account_subscription label: %q", got)
	}
	if got := inviteTypeLabel(langRU, "roster_only"); got != "Только контакт" {
		t.Fatalf("unexpected roster_only label: %q", got)
	}
	if got := inviteTypeLabel(langEN, "future_type"); got != "future type" {
		t.Fatalf("unexpected fallback invite type label: %q", got)
	}
}

func TestBasePath(t *testing.T) {
	cases := map[string]string{
		"":            "",
		"/":           "",
		"xmpp-admin":  "/xmpp-admin",
		"/xmpp-admin": "/xmpp-admin",
		"/xmpp-admin/": "/xmpp-admin",
	}
	for in, want := range cases {
		if got := normalizeBasePath(in); got != want {
			t.Errorf("normalizeBasePath(%q) = %q, want %q", in, got, want)
		}
	}
	if got := joinBasePath("/xmpp-admin", "/admin"); got != "/xmpp-admin/admin" {
		t.Fatalf("unexpected joined path: %q", got)
	}
}

func TestSameOriginRejectsCrossSitePost(t *testing.T) {
	req := httptest.NewRequest("POST", "http://admin.example/x", nil)
	req.Host = "admin.example"
	req.Header.Set("X-Forwarded-Proto", "https")
	req.Header.Set("Origin", "https://evil.example")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	if sameOrigin(req) {
		t.Fatal("cross-site request must be rejected")
	}
}

func TestSameOriginAcceptsProxyHTTPS(t *testing.T) {
	req := httptest.NewRequest("POST", "http://admin.example/x", nil)
	req.Host = "admin.example"
	req.Header.Set("X-Forwarded-Proto", "https")
	req.Header.Set("Origin", "https://admin.example")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	if !sameOrigin(req) {
		t.Fatal("same-origin proxy request should be accepted")
	}
}


func TestTemplatesParse(t *testing.T) {
	_, err := template.New("root").Funcs(template.FuncMap{
		"tr":              tr,
		"join":            strings.Join,
		"p":               func(path string) string { return path },
		"safeURL":         safeExternalURL,
		"formatTimestamp": formatTimestamp,
		"inviteTypeLabel": inviteTypeLabel,
	}).Parse(adminTemplate + opsTemplate)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTemplateTranslationKeysExist(t *testing.T) {
	re := regexp.MustCompile(`tr\s+(?:\$?\.Lang)\s+"([^"]+)"`)
	matches := re.FindAllStringSubmatch(adminTemplate+opsTemplate, -1)
	for _, match := range matches {
		key := match[1]
		if _, ok := messages[langEN][key]; !ok {
			t.Errorf("template translation key %q is missing from EN catalog", key)
		}
		if _, ok := messages[langRU][key]; !ok {
			t.Errorf("template translation key %q is missing from RU catalog", key)
		}
	}
}


func TestSafeExternalURL(t *testing.T) {
	accepted := []string{
		"https://example.org/invites/token",
		"http://127.0.0.1:5280/invites/token",
		"xmpp:example.org?register;preauth=token",
	}
	for _, raw := range accepted {
		if got := safeExternalURL(raw); got == "" {
			t.Errorf("expected URL to be accepted: %q", raw)
		}
	}

	rejected := []string{
		"javascript:alert(1)",
		"file:///etc/passwd",
		"data:text/html,x",
		"//evil.example/path",
		"https:///missing-host",
	}
	for _, raw := range rejected {
		if got := safeExternalURL(raw); got != "" {
			t.Errorf("expected URL to be rejected: %q -> %q", raw, got)
		}
	}
}


func TestUsernameValidationDefersInternationalNormalizationToEjabberd(t *testing.T) {
	valid := []string{"alice", "álîçé", "юзер", "user.name", "name_1"}
	for _, value := range valid {
		if !validUsername(value) {
			t.Errorf("expected username %q to be accepted for server-side validation", value)
		}
	}

	invalid := []string{"", "user@example.org", "user/resource", "two words", "line\nbreak"}
	for _, value := range invalid {
		if validUsername(value) {
			t.Errorf("expected username %q to be rejected", value)
		}
	}
}


func TestSameOriginAllowsMissingBrowserMetadataForCSRFFallback(t *testing.T) {
	req := httptest.NewRequest("POST", "http://admin.example/x", nil)
	req.Host = "admin.example"
	req.Header.Set("X-Forwarded-Proto", "https")
	if !sameOrigin(req) {
		t.Fatal("missing browser metadata must fall back to CSRF token validation")
	}
}

func TestRequireAdminRejectsPOSTWithoutCSRFToken(t *testing.T) {
	app := &App{
		cfg: Config{
			AdminUser:     "admin",
			AdminPassword: "a-very-long-admin-password",
		},
		csrfToken: "test-csrf-token",
	}
	req := httptest.NewRequest(http.MethodPost, "https://admin.example/admin/invites/create", strings.NewReader(""))
	req.Host = "admin.example"
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("admin", "a-very-long-admin-password")

	rec := httptest.NewRecorder()
	if app.requireAdmin(rec, req) {
		t.Fatal("POST without CSRF token must be rejected")
	}
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusForbidden)
	}
}

func TestRequireAdminAcceptsPOSTWithCSRFTokenWithoutBrowserMetadata(t *testing.T) {
	app := &App{
		cfg: Config{
			AdminUser:     "admin",
			AdminPassword: "a-very-long-admin-password",
		},
		csrfToken: "test-csrf-token",
	}
	form := url.Values{"csrf_token": {"test-csrf-token"}}
	req := httptest.NewRequest(http.MethodPost, "https://admin.example/admin/invites/create", strings.NewReader(form.Encode()))
	req.Host = "admin.example"
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("admin", "a-very-long-admin-password")

	rec := httptest.NewRecorder()
	if !app.requireAdmin(rec, req) {
		t.Fatalf("valid CSRF token should allow POST without browser metadata; status=%d", rec.Code)
	}
}

func TestSameOriginAcceptsFetchMetadataFallback(t *testing.T) {
	req := httptest.NewRequest("POST", "http://admin.example/x", nil)
	req.Host = "admin.example"
	req.Header.Set("X-Forwarded-Proto", "https")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	if !sameOrigin(req) {
		t.Fatal("same-origin Fetch Metadata should be accepted")
	}
}

func TestSameOriginAcceptsPrivacyRedactedMetadata(t *testing.T) {
	req := httptest.NewRequest("POST", "http://admin.example/x", nil)
	req.Host = "admin.example"
	req.Header.Set("X-Forwarded-Proto", "https")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Origin", "null")
	if !sameOrigin(req) {
		t.Fatal("privacy-redacted metadata must fall back to CSRF validation")
	}
}

func TestSameOriginRejectsExplicitForeignOrigin(t *testing.T) {
	req := httptest.NewRequest("POST", "http://admin.example/x", nil)
	req.Host = "admin.example"
	req.Header.Set("X-Forwarded-Proto", "https")
	req.Header.Set("Origin", "https://evil.example")
	if sameOrigin(req) {
		t.Fatal("explicit foreign origin must be rejected")
	}
}



func TestCreateInviteRedirectDoesNotLeakToken(t *testing.T) {
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/generate_invite" {
			t.Fatalf("unexpected API path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"invite_uri":"xmpp:example.org?register;preauth=secret-token","landing_page":"https://example.org/invites/secret-token"}`))
	}))
	defer api.Close()

	app := &App{
		cfg: Config{
			Domain:             "example.org",
			AdminUser:          "admin",
			AdminPassword:      "a-very-long-admin-password",
			EjabberdConfigPath: filepath.Join(t.TempDir(), "missing.yml"),
			Server:             EjabberdConfigSnapshot{InvitesEnabled: true},
		},
		ejabberd: &EjabberdClient{
			baseURL: api.URL,
			client:  api.Client(),
		},
		csrfToken: "test-csrf-token",
	}

	form := url.Values{"csrf_token": {"test-csrf-token"}}
	req := httptest.NewRequest(http.MethodPost, "https://admin.example/admin/invites/create", strings.NewReader(form.Encode()))
	req.Host = "admin.example"
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://admin.example")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.SetBasicAuth("admin", "a-very-long-admin-password")

	rec := httptest.NewRecorder()
	app.adminCreateInvite(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusSeeOther)
	}
	location := rec.Header().Get("Location")
	if location != "/admin?created=1" {
		t.Fatalf("unexpected redirect location: %q", location)
	}
	if strings.Contains(location, "secret-token") {
		t.Fatal("invite token leaked into redirect URL")
	}
}
