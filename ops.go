package main

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type UserRow struct {
	Username     string
	JID          string
	SessionCount int
	SessionState string
}

type SessionRow struct {
	JID      string
	Username string
	Resource string
}

type RoomRow struct {
	JID     string
	Name    string
	Service string
}

type HealthCheck struct {
	Name      string
	State     string
	StateKey  string
	Summary   string
	Technical string
}

type OpsPageData struct {
	Lang                 string
	Active               string
	CurrentPath          string
	TitleKey             string
	SubtitleKey          string
	Domain               string
	APIError             string
	APIErrorTitle        string
	Warnings             []string
	Users                []UserRow
	Sessions             []SessionRow
	Rooms                []RoomRow
	Health               []HealthCheck
	RegisteredMetric     string
	OnlineSessionsMetric string
	OnlineRoomsMetric    string
	APILatency           string
}

func (a *App) adminUsers(w http.ResponseWriter, r *http.Request) {
	if !a.requireAdminGET(w, r, "/admin/users") {
		return
	}
	lang := requestLanguage(r)
	data := OpsPageData{
		Lang:        lang,
		Active:      "users",
		CurrentPath: "/admin/users",
		TitleKey:    "users",
		SubtitleKey: "users_subtitle",
		Domain:      a.cfg.Domain,
	}

	users, err := a.ejabberd.RegisteredUsers(r.Context(), a.cfg.Domain)
	if err != nil {
		data.APIErrorTitle, data.APIError = apiProblem(lang, "users", err)
		a.renderOps(w, data)
		return
	}
	sort.Strings(users)

	sessionCounts := map[string]int{}
	sessions, sessionsErr := a.ejabberd.ConnectedUsers(r.Context())
	if sessionsErr != nil {
		data.Warnings = append(data.Warnings, tr(lang, "sessions_lookup_unavailable"))
	} else {
		for _, session := range sessionsForHost(sessions, a.cfg.Domain) {
			if session.Username != "" {
				sessionCounts[session.Username]++
			}
		}
	}

	data.Users = make([]UserRow, 0, len(users))
	for _, user := range users {
		count, known := sessionCounts[user]
		state := "unknown"
		if sessionsErr == nil {
			state = "offline"
			if known && count > 0 {
				state = "online"
			}
		}
		data.Users = append(data.Users, UserRow{
			Username:     user,
			JID:          user + "@" + a.cfg.Domain,
			SessionCount: count,
			SessionState: state,
		})
	}
	a.renderOps(w, data)
}

func (a *App) adminSessions(w http.ResponseWriter, r *http.Request) {
	if !a.requireAdminGET(w, r, "/admin/sessions") {
		return
	}
	lang := requestLanguage(r)
	data := OpsPageData{
		Lang:        lang,
		Active:      "sessions",
		CurrentPath: "/admin/sessions",
		TitleKey:    "connections_title",
		SubtitleKey: "sessions_subtitle",
		Domain:      a.cfg.Domain,
	}

	sessions, err := a.ejabberd.ConnectedUsers(r.Context())
	if err != nil {
		data.APIErrorTitle, data.APIError = apiProblem(lang, "sessions", err)
		a.renderOps(w, data)
		return
	}
	data.Sessions = sessionsForHost(sessions, a.cfg.Domain)
	a.renderOps(w, data)
}

func (a *App) adminRooms(w http.ResponseWriter, r *http.Request) {
	if !a.requireAdminGET(w, r, "/admin/rooms") {
		return
	}
	lang := requestLanguage(r)
	data := OpsPageData{
		Lang:        lang,
		Active:      "rooms",
		CurrentPath: "/admin/rooms",
		TitleKey:    "rooms",
		SubtitleKey: "rooms_subtitle",
		Domain:      a.cfg.Domain,
	}

	rooms, err := a.ejabberd.MUCOnlineRooms(r.Context(), "global")
	if err != nil {
		data.APIErrorTitle, data.APIError = apiProblem(lang, "rooms", err)
		a.renderOps(w, data)
		return
	}
	sort.Strings(rooms)
	for _, jid := range rooms {
		name, service, _ := splitFullJID(jid)
		data.Rooms = append(data.Rooms, RoomRow{JID: jid, Name: name, Service: service})
	}
	a.renderOps(w, data)
}

func (a *App) adminHealth(w http.ResponseWriter, r *http.Request) {
	if !a.requireAdminGET(w, r, "/admin/health") {
		return
	}
	lang := requestLanguage(r)
	data := OpsPageData{
		Lang:                 lang,
		Active:               "health",
		CurrentPath:          "/admin/health",
		TitleKey:             "infrastructure_health",
		SubtitleKey:          "health_subtitle",
		Domain:               a.cfg.Domain,
		RegisteredMetric:     "—",
		OnlineSessionsMetric: "—",
		OnlineRoomsMetric:    "—",
		APILatency:           "—",
	}

	data.Health = append(data.Health, HealthCheck{
		Name:     tr(lang, "xmpp_admin_service"),
		State:    "ok",
		StateKey: "healthy",
		Summary:  tr(lang, "xmpp_admin_running"),
	})

	snapshot, configErr := loadEjabberdConfigForHost(a.cfg.EjabberdConfigPath, a.cfg.Domain)
	if configErr != nil {
		data.Health = append(data.Health, HealthCheck{
			Name:      tr(lang, "configuration"),
			State:     "error",
			StateKey:  "unavailable",
			Summary:   tr(lang, "health_config_unavailable"),
			Technical: configErr.Error(),
		})
	} else {
		data.Health = append(data.Health, HealthCheck{
			Name:      tr(lang, "configuration"),
			State:     "ok",
			StateKey:  "healthy",
			Summary:   tr(lang, "read_from_config"),
			Technical: snapshot.Path,
		})
		inviteState := "warn"
		inviteKey := "disabled"
		if snapshot.InvitesEnabled {
			inviteState = "ok"
			inviteKey = "enabled"
		}
		data.Health = append(data.Health, HealthCheck{
			Name:     tr(lang, "invitations"),
			State:    inviteState,
			StateKey: inviteKey,
			Summary:  tr(lang, "read_from_config"),
		})
	}

	start := time.Now()
	statusErr := a.ejabberd.Status(r.Context())
	statusLatency := time.Since(start).Round(time.Millisecond).String()
	if statusErr != nil {
		summary, technical := apiProblem(lang, "api", statusErr)
		data.Health = append(data.Health, HealthCheck{
			Name:      tr(lang, "ejabberd_api"),
			State:     "error",
			StateKey:  "unavailable",
			Summary:   summary,
			Technical: technical,
		})
	} else {
		data.APILatency = statusLatency
		data.Health = append(data.Health, HealthCheck{
			Name:     tr(lang, "ejabberd_api"),
			State:    "ok",
			StateKey: "healthy",
			Summary:  tr(lang, "status_command_ok"),
		})
	}

	users, usersErr := a.ejabberd.RegisteredUsers(r.Context(), a.cfg.Domain)
	if usersErr != nil {
		summary, technical := apiProblem(lang, "users", usersErr)
		data.Health = append(data.Health, HealthCheck{
			Name:      tr(lang, "users"),
			State:     "error",
			StateKey:  "unavailable",
			Summary:   summary,
			Technical: technical,
		})
	} else {
		data.RegisteredMetric = strconv.Itoa(len(users))
		data.Health = append(data.Health, HealthCheck{
			Name:     tr(lang, "users"),
			State:    "ok",
			StateKey: "healthy",
			Summary:  tr(lang, "read_access_ok"),
		})
	}

	sessions, sessionsErr := a.ejabberd.ConnectedUsers(r.Context())
	if sessionsErr != nil {
		summary, technical := apiProblem(lang, "sessions", sessionsErr)
		data.Health = append(data.Health, HealthCheck{
			Name:      tr(lang, "connections_title"),
			State:     "error",
			StateKey:  "unavailable",
			Summary:   summary,
			Technical: technical,
		})
	} else {
		data.OnlineSessionsMetric = strconv.Itoa(len(sessionsForHost(sessions, a.cfg.Domain)))
		data.Health = append(data.Health, HealthCheck{
			Name:     tr(lang, "connections_title"),
			State:    "ok",
			StateKey: "healthy",
			Summary:  tr(lang, "read_access_ok"),
		})
	}

	rooms, roomsErr := a.ejabberd.MUCOnlineRooms(r.Context(), "global")
	if roomsErr != nil {
		summary, technical := apiProblem(lang, "rooms", roomsErr)
		data.Health = append(data.Health, HealthCheck{
			Name:      tr(lang, "rooms"),
			State:     "warn",
			StateKey:  "unavailable",
			Summary:   summary,
			Technical: technical,
		})
	} else {
		data.OnlineRoomsMetric = strconv.Itoa(len(rooms))
		data.Health = append(data.Health, HealthCheck{
			Name:     tr(lang, "rooms"),
			State:    "ok",
			StateKey: "healthy",
			Summary:  tr(lang, "read_access_ok"),
		})
	}

	sort.SliceStable(data.Health, func(i, j int) bool {
		return healthStateRank(data.Health[i].State) < healthStateRank(data.Health[j].State)
	})
	a.renderOps(w, data)
}

func healthStateRank(state string) int {
	switch state {
	case "error":
		return 0
	case "warn":
		return 1
	default:
		return 2
	}
}

func apiProblem(lang, area string, err error) (string, string) {
	technical := ""
	if err != nil {
		technical = err.Error()
	}
	forbidden := strings.Contains(technical, "HTTP 403")
	key := ""
	switch area {
	case "api":
		if forbidden {
			key = "health_api_forbidden"
		} else {
			key = "health_api_unavailable"
		}
	case "users":
		if forbidden {
			key = "health_users_forbidden"
		} else {
			key = "health_users_unavailable"
		}
	case "sessions":
		if forbidden {
			key = "health_sessions_forbidden"
		} else {
			key = "health_sessions_unavailable"
		}
	case "rooms":
		if forbidden {
			key = "health_rooms_forbidden"
		} else {
			key = "health_rooms_unavailable"
		}
	default:
		key = "data_unavailable"
	}
	return tr(lang, key), technical
}

func (a *App) requireAdminGET(w http.ResponseWriter, r *http.Request, path string) bool {
	if !a.requireAdmin(w, r) {
		return false
	}
	if r.URL.Path != path {
		http.NotFound(w, r)
		return false
	}
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "GET required", http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func (a *App) renderOps(w http.ResponseWriter, data OpsPageData) {
	if err := a.templates.ExecuteTemplate(w, "ops", data); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func sessionsForHost(rawSessions []string, domain string) []SessionRow {
	rows := make([]SessionRow, 0, len(rawSessions))
	for _, raw := range rawSessions {
		user, host, resource := splitFullJID(raw)
		if !strings.EqualFold(host, domain) {
			continue
		}
		rows = append(rows, SessionRow{
			JID:      raw,
			Username: user,
			Resource: resource,
		})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].JID < rows[j].JID })
	return rows
}

func splitFullJID(raw string) (user, host, resource string) {
	raw = strings.TrimSpace(raw)
	bare := raw
	if slash := strings.IndexByte(raw, '/'); slash >= 0 {
		bare = raw[:slash]
		resource = raw[slash+1:]
	}
	at := strings.LastIndexByte(bare, '@')
	if at < 0 {
		return bare, "", resource
	}
	return bare[:at], bare[at+1:], resource
}
