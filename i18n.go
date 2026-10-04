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
