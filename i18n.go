package main

import (
	"net/http"
	"strings"
	"time"
)

const (
	langEN = "en"
	langRU = "ru"
)

var messages = map[string]map[string]string{
	langEN: {
		"invitations":            "Invitations",
		"invitations_subtitle":   "Native ejabberd invitations for %s",
		"api_error":              "Could not read invitations from ejabberd",
		"invite_created":         "Invite created",
		"native_created_note":    "The invitation is stored and enforced by ejabberd.",
		"native_invites_note":    "Invitations are read directly from ejabberd. Invites created by registered users in XMPP clients appear here automatically.",
		"new_invite":             "New invite",
		"username_optional":      "Preselected username",
		"username_optional_hint": "Leave empty to let the recipient choose a username.",
		"create_invite":          "Create invite",
		"policy_from_ejabberd":   "Lifetime, limits and access policy are taken from the active ejabberd configuration.",
		"server_config":          "Server configuration",
		"domain":                 "Domain",
		"enabled":                "enabled",
		"disabled":               "disabled",
		"access_rule":            "Invite access rule",
		"max_invites":            "Per-user invite limit",
		"token_ttl":              "Invite lifetime",
		"seconds":                "seconds",
		"landing_page":           "Landing page",
		"templates_dir":          "Templates",
		"site_name":              "Site name",
		"db_type":                "Database",
		"webchat_url":            "Webchat",
		"allow_modules":          "Allowed registration modules",
		"ejabberd_default":       "ejabberd default",
		"config_path":            "Configuration file",
		"created":                "Created",
		"inviter":                "Inviter",
		"account_name":           "Reserved username",
		"expires":                "Expires",
		"type":                   "Type",
		"status":                 "Status",
		"server_generated":       "Server / administrator",
		"active":                 "active",
		"expired":                "expired",
		"open":                   "Open",
		"revoke":                 "Revoke",
		"no_invites":             "No invitations",
		"theme_dark":             "Use dark theme",
		"theme_light":            "Use light theme",
		"mod_invites_disabled":   "mod_invites is not enabled in the loaded ejabberd configuration.",
		"invalid_username":       "Invalid username.",
		"invite_type_account_subscription": "Account + contact",
		"invite_type_account_only":         "Account only",
		"invite_type_roster_only":          "Contact only",
		"admin_sections":                    "Administration sections",
		"nav_invites":                       "Invites",
		"nav_users":                         "Users",
		"nav_sessions":                      "Sessions",
		"nav_rooms":                         "Rooms",
		"nav_health":                        "Health",
		"users":                             "Users",
		"sessions":                          "Sessions",
		"rooms":                             "Rooms",
		"health":                            "Health",
		"users_subtitle":                    "Registered accounts for %s",
		"sessions_subtitle":                 "Active sessions for %s",
		"rooms_subtitle":                    "Server-wide online MUC rooms reported by ejabberd (selected vhost: %s)",
		"health_subtitle":                   "Operational checks for %s",
		"data_unavailable":                  "Could not load data from ejabberd",
		"sessions_lookup_unavailable":       "Session state is temporarily unavailable; registered accounts are still shown.",
		"registered_accounts":               "Registered accounts",
		"account":                           "Account",
		"jid":                               "JID",
		"connection":                        "Connection",
		"connected":                         "connected",
		"not_connected":                     "not connected",
		"online":                            "online",
		"offline":                           "offline",
		"no_users":                          "No registered accounts",
		"active_sessions":                   "Active sessions",
		"resource":                          "Resource",
		"full_jid":                          "Full JID",
		"no_sessions":                       "No active sessions",
		"online_rooms":                      "Online rooms",
		"room":                              "Room",
		"service":                           "Service",
		"no_rooms":                          "No online rooms",
		"infrastructure_health":             "Infrastructure health",
		"registered_users_metric":           "Registered users",
		"online_sessions_metric":            "Online sessions",
		"online_rooms_metric":               "Online rooms (server)",
		"api_latency":                       "API latency",
		"ejabberd_api_latency":              "ejabberd API latency",
		"checks":                            "Checks",
		"component":                         "Component",
		"detail":                            "Detail",
		"healthy":                           "healthy",
		"unavailable":                       "unavailable",
		"xmpp_admin_service":                "XMPP Admin",
		"xmpp_admin_running":                "Web service is running.",
		"configuration":                     "Configuration",
		"read_from_config":                  "Read from ejabberd configuration.",
		"ejabberd_api":                      "ejabberd API",
		"status_command_ok":                 "Status command completed successfully.",
		"read_access_ok":                    "Read access is working.",
		"rooms_command_optional":            "Room listing is unavailable. Ensure mod_muc_admin and the muc_online_rooms API permission are enabled.",
	},
	langRU: {
		"invitations":            "Приглашения",
		"invitations_subtitle":   "Нативные инвайты ejabberd для %s",
		"api_error":              "Не удалось получить инвайты из ejabberd",
		"invite_created":         "Инвайт создан",
		"native_created_note":    "Инвайт хранится и проверяется самим ejabberd.",
		"native_invites_note":    "Инвайты читаются напрямую из ejabberd. Инвайты, созданные зарегистрированными пользователями в XMPP-клиентах, автоматически появляются здесь.",
		"new_invite":             "Новый инвайт",
		"username_optional":      "Заранее заданный username",
		"username_optional_hint": "Оставьте пустым, чтобы получатель сам выбрал username.",
		"create_invite":          "Создать инвайт",
		"policy_from_ejabberd":   "Срок действия, лимиты и права берутся из активной конфигурации ejabberd.",
		"server_config":          "Конфигурация сервера",
		"domain":                 "Домен",
		"enabled":                "включён",
		"disabled":               "выключен",
		"access_rule":            "Правило доступа к инвайтам",
		"max_invites":            "Лимит инвайтов на пользователя",
		"token_ttl":              "Срок действия инвайта",
		"seconds":                "секунд",
		"landing_page":           "Страница приглашения",
		"templates_dir":          "Шаблоны",
		"site_name":              "Название сайта",
		"db_type":                "База данных",
		"webchat_url":            "Webchat",
		"allow_modules":          "Разрешённые модули регистрации",
		"ejabberd_default":       "стандартные ejabberd",
		"config_path":            "Файл конфигурации",
		"created":                "Создан",
		"inviter":                "Кто создал",
		"account_name":           "Зарезервированный username",
		"expires":                "Истекает",
		"type":                   "Тип",
		"status":                 "Статус",
		"server_generated":       "Сервер / администратор",
		"active":                 "активен",
		"expired":                "истёк",
		"open":                   "Открыть",
		"revoke":                 "Отозвать",
		"no_invites":             "Инвайтов нет",
		"theme_dark":             "Включить тёмную тему",
		"theme_light":            "Включить светлую тему",
		"mod_invites_disabled":   "mod_invites не включён в загруженной конфигурации ejabberd.",
		"invalid_username":       "Некорректный username.",
		"invite_type_account_subscription": "Аккаунт + контакт",
		"invite_type_account_only":         "Только аккаунт",
		"invite_type_roster_only":          "Только контакт",
		"admin_sections":                    "Разделы администрирования",
		"nav_invites":                       "Инвайты",
		"nav_users":                         "Люди",
		"nav_sessions":                      "Сессии",
		"nav_rooms":                         "Комнаты",
		"nav_health":                        "Статус",
		"users":                             "Пользователи",
		"sessions":                          "Сессии",
		"rooms":                             "Комнаты",
		"health":                            "Состояние",
		"users_subtitle":                    "Зарегистрированные аккаунты для %s",
		"sessions_subtitle":                 "Активные сессии для %s",
		"rooms_subtitle":                    "Онлайн MUC-комнаты на сервере (выбранный vhost: %s)",
		"health_subtitle":                   "Проверки инфраструктуры для %s",
		"data_unavailable":                  "Не удалось получить данные из ejabberd",
		"sessions_lookup_unavailable":       "Состояние сессий временно недоступно; список зарегистрированных аккаунтов всё равно показан.",
		"registered_accounts":               "Зарегистрированные аккаунты",
		"account":                           "Аккаунт",
		"jid":                               "JID",
		"connection":                        "Подключение",
		"connected":                         "подключён",
		"not_connected":                     "не подключён",
		"online":                            "в сети",
		"offline":                           "не в сети",
		"no_users":                          "Зарегистрированных аккаунтов нет",
		"active_sessions":                   "Активные сессии",
		"resource":                          "Ресурс",
		"full_jid":                          "Полный JID",
		"no_sessions":                       "Активных сессий нет",
		"online_rooms":                      "Онлайн-комнаты",
		"room":                              "Комната",
		"service":                           "Сервис",
		"no_rooms":                          "Онлайн-комнат нет",
		"infrastructure_health":             "Состояние инфраструктуры",
		"registered_users_metric":           "Пользователей",
		"online_sessions_metric":            "Сессий онлайн",
		"online_rooms_metric":               "Комнат онлайн (сервер)",
		"api_latency":                       "Задержка API",
		"ejabberd_api_latency":              "Задержка API ejabberd",
		"checks":                            "Проверки",
		"component":                         "Компонент",
		"detail":                            "Подробности",
		"healthy":                           "работает",
		"unavailable":                       "недоступно",
		"xmpp_admin_service":                "XMPP Admin",
		"xmpp_admin_running":                "Веб-сервис работает.",
		"configuration":                     "Конфигурация",
		"read_from_config":                  "Прочитано из конфигурации ejabberd.",
		"ejabberd_api":                      "API ejabberd",
		"status_command_ok":                 "Команда status выполнена успешно.",
		"read_access_ok":                    "Доступ на чтение работает.",
		"rooms_command_optional":            "Список комнат недоступен. Проверьте mod_muc_admin и разрешение API muc_online_rooms.",
	},
}

func tr(lang, key string) string {
	if lang != langRU {
		lang = langEN
	}
	if value, ok := messages[lang][key]; ok {
		return value
	}
	if value, ok := messages[langEN][key]; ok {
		return value
	}
	return key
}

func inviteTypeLabel(lang, raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "account_subscription":
		return tr(lang, "invite_type_account_subscription")
	case "account_only":
		return tr(lang, "invite_type_account_only")
	case "roster_only":
		return tr(lang, "invite_type_roster_only")
	default:
		raw = strings.TrimSpace(raw)
		if raw == "" {
			return "—"
		}
		return strings.ReplaceAll(raw, "_", " ")
	}
}

func formatTimestamp(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "—"
	}
	if parsed, err := time.Parse(time.RFC3339Nano, raw); err == nil {
		return parsed.UTC().Format("2006-01-02 15:04") + " UTC"
	}
	return raw
}

func requestLanguage(r *http.Request) string {
	if cookie, err := r.Cookie("xmpp_admin_lang"); err == nil {
		if cookie.Value == langRU || cookie.Value == langEN {
			return cookie.Value
		}
	}
	for _, part := range strings.Split(r.Header.Get("Accept-Language"), ",") {
		tag := strings.ToLower(strings.TrimSpace(strings.SplitN(part, ";", 2)[0]))
		if tag == "ru" || strings.HasPrefix(tag, "ru-") {
			return langRU
		}
		if tag == "en" || strings.HasPrefix(tag, "en-") {
			return langEN
		}
	}
	return langEN
}

func setLanguageCookie(w http.ResponseWriter, r *http.Request, lang string) {
	if lang != langRU {
		lang = langEN
	}
	secure := r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https")
	http.SetCookie(w, &http.Cookie{
		Name:     "xmpp_admin_lang",
		Value:    lang,
		Path:     "/",
		MaxAge:   int((365 * 24 * time.Hour).Seconds()),
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}

func safeNext(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" || !strings.HasPrefix(raw, "/") || strings.HasPrefix(raw, "//") {
		return "/admin"
	}
	return raw
}
