package main

import (
	"fmt"
	"net/http"
	"strconv"
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
		"invitations_subtitle":   "Create, share and revoke invitations for %s",
		"api_error":              "Could not load invitations",
		"invite_created":         "Invitation ready",
		"native_created_note":    "The newest invitation is highlighted below.",
		"native_invites_note":    "Invitations created here or from XMPP clients appear in the same list.",
		"new_invite":             "Create invitation",
		"username_optional":      "Account name (optional)",
		"username_optional_hint": "Leave empty and the recipient can choose the account name.",
		"create_invite":          "Create invitation",
		"policy_from_ejabberd":   "Current limits and lifetime are shown in Invitation settings.",
		"server_config":          "Invitation settings",
		"domain":                 "Domain",
		"enabled":                "enabled",
		"disabled":               "disabled",
		"access_rule":            "Who can create",
		"max_invites":            "Per-user limit",
		"token_ttl":              "Valid for",
		"seconds":                "seconds",
		"landing_page":           "Invite page",
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
		"nav_sessions":                      "Connections",
		"nav_rooms":                         "Rooms",
		"nav_health":                        "Health",
		"connections_title":                 "Connections",
		"search":                            "Search",
		"search_users_placeholder":          "Search accounts or JIDs",
		"search_sessions_placeholder":       "Search account, JID or resource",
		"search_rooms_placeholder":          "Search rooms or services",
		"filter_all":                        "All",
		"filter_active":                     "Active",
		"filter_expired":                    "Expired",
		"copy_link":                         "Copy link",
		"copied":                            "Copied",
		"invite_ready":                      "New invitation",
		"recipient_chooses_name":            "Recipient chooses the account name",
		"created_by":                        "Created by",
		"created_when":                      "Created",
		"expires_when":                      "Expires",
		"no_invites_title":                  "No invitations yet",
		"no_invites_body":                   "Create an invitation and it will appear here.",
		"no_users_title":                    "No accounts yet",
		"no_users_body":                     "Registered accounts will appear here.",
		"no_sessions_title":                 "No active connections",
		"no_sessions_body":                  "Device connections will appear here automatically.",
		"no_rooms_title":                    "No active rooms",
		"no_rooms_body":                     "Rooms will appear here when they become active.",
		"nothing_found":                     "Nothing found",
		"try_different_search":               "Try a different search.",
		"no_filtered_invites":                "No invitations with this status.",
		"technical_details":                 "Technical details",
		"all_users":                         "All users",
		"no_one":                            "Nobody",
		"no_limit":                          "No limit",
		"automatic":                         "Automatic",
		"not_used":                          "Not used",
		"configured":                        "Configured",
		"server_policy":                     "Restricted by server policy",
		"health_config_unavailable":         "Could not read the server configuration.",
		"health_api_forbidden":              "XMPP Admin does not have permission to read ejabberd status.",
		"health_api_unavailable":            "The ejabberd API is not responding.",
		"health_users_forbidden":            "XMPP Admin does not have permission to read users.",
		"health_users_unavailable":          "Could not load the user list.",
		"health_sessions_forbidden":         "XMPP Admin does not have permission to read connections.",
		"health_sessions_unavailable":       "Could not load active connections.",
		"health_rooms_forbidden":            "XMPP Admin does not have permission to read rooms.",
		"health_rooms_unavailable":          "Could not load active rooms.",
		"users":                             "Users",
		"sessions":                          "Sessions",
		"rooms":                             "Rooms",
		"health":                            "Health",
		"users_subtitle":                    "Find an account and check its connection state on %s",
		"sessions_subtitle":                 "See active device connections to %s",
		"rooms_subtitle":                    "See active group chats on the server for %s",
		"health_subtitle":                   "See what is working and what needs attention for %s",
		"data_unavailable":                  "Could not load data from ejabberd",
		"sessions_lookup_unavailable":       "Session state is temporarily unavailable; registered accounts are still shown.",
		"registered_accounts":               "Accounts",
		"account":                           "Account",
		"jid":                               "JID",
		"connection":                        "Connection",
		"connected":                         "connected",
		"not_connected":                     "not connected",
		"online":                            "online",
		"offline":                           "offline",
		"no_users":                          "No registered accounts",
		"active_sessions":                   "Active connections",
		"resource":                          "Device / resource",
		"full_jid":                          "Full JID",
		"no_sessions":                       "No active sessions",
		"online_rooms":                      "Active rooms",
		"room":                              "Room",
		"service":                           "Service",
		"no_rooms":                          "No online rooms",
		"infrastructure_health":             "Server health",
		"registered_users_metric":           "Accounts",
		"online_sessions_metric":            "Active connections",
		"online_rooms_metric":               "Active rooms (server)",
		"api_latency":                       "API latency",
		"ejabberd_api_latency":              "ejabberd API latency",
		"checks":                            "Services and checks",
		"component":                         "Component",
		"detail":                            "Detail",
		"healthy":                           "healthy",
		"unavailable":                       "unavailable",
		"xmpp_admin_service":                "XMPP Admin",
		"xmpp_admin_running":                "Admin panel is responding.",
		"configuration":                     "Server settings",
		"read_from_config":                  "Settings loaded.",
		"ejabberd_api":                      "ejabberd connection",
		"status_command_ok":                 "Connection to ejabberd is working.",
		"read_access_ok":                    "Data is available.",
		"rooms_command_optional":            "Room listing is unavailable. Ensure mod_muc_admin and the muc_online_rooms API permission are enabled.",
	},
	langRU: {
		"invitations":            "Приглашения",
		"invitations_subtitle":   "Создавайте, отправляйте и отзывайте приглашения для %s",
		"api_error":              "Не удалось загрузить приглашения",
		"invite_created":         "Приглашение готово",
		"native_created_note":    "Новое приглашение выделено первым в списке.",
		"native_invites_note":    "Приглашения, созданные здесь или в XMPP-клиентах, появляются в одном списке.",
		"new_invite":             "Создать приглашение",
		"username_optional":      "Имя аккаунта (необязательно)",
		"username_optional_hint": "Оставьте пустым — получатель сам выберет имя аккаунта.",
		"create_invite":          "Создать приглашение",
		"policy_from_ejabberd":   "Текущие лимиты и срок действия показаны в настройках приглашений.",
		"server_config":          "Настройки приглашений",
		"domain":                 "Домен",
		"enabled":                "включён",
		"disabled":               "выключен",
		"access_rule":            "Кто может создавать",
		"max_invites":            "Лимит на пользователя",
		"token_ttl":              "Срок действия",
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
		"nav_invites":                       "Приглашения",
		"nav_users":                         "Пользователи",
		"nav_sessions":                      "Подключения",
		"nav_rooms":                         "Комнаты",
		"nav_health":                        "Состояние",
		"connections_title":                 "Подключения",
		"search":                            "Поиск",
		"search_users_placeholder":          "Поиск по аккаунту или JID",
		"search_sessions_placeholder":       "Поиск по аккаунту, JID или ресурсу",
		"search_rooms_placeholder":          "Поиск по комнате или сервису",
		"filter_all":                        "Все",
		"filter_active":                     "Активные",
		"filter_expired":                    "Истёкшие",
		"copy_link":                         "Скопировать ссылку",
		"copied":                            "Скопировано",
		"invite_ready":                      "Новое приглашение",
		"recipient_chooses_name":            "Имя аккаунта выберет получатель",
		"created_by":                        "Создал",
		"created_when":                      "Создано",
		"expires_when":                      "Истекает",
		"no_invites_title":                  "Приглашений пока нет",
		"no_invites_body":                   "Создайте приглашение — оно появится здесь.",
		"no_users_title":                    "Пользователей пока нет",
		"no_users_body":                     "Зарегистрированные аккаунты появятся здесь.",
		"no_sessions_title":                 "Сейчас никто не подключён",
		"no_sessions_body":                  "Активные подключения устройств появятся здесь автоматически.",
		"no_rooms_title":                    "Сейчас нет активных комнат",
		"no_rooms_body":                     "Комнаты появятся здесь, когда станут активны.",
		"nothing_found":                     "Ничего не найдено",
		"try_different_search":               "Попробуйте изменить запрос.",
		"no_filtered_invites":                "Приглашений с таким статусом нет.",
		"technical_details":                 "Технические сведения",
		"all_users":                         "Все пользователи",
		"no_one":                            "Никто",
		"no_limit":                          "Без лимита",
		"automatic":                         "Автоматически",
		"not_used":                          "Не используется",
		"configured":                        "Настроена",
		"server_policy":                     "Ограничено правилами сервера",
		"health_config_unavailable":         "Не удалось прочитать конфигурацию сервера.",
		"health_api_forbidden":              "У XMPP Admin нет права читать состояние ejabberd.",
		"health_api_unavailable":            "API ejabberd не отвечает.",
		"health_users_forbidden":            "У XMPP Admin нет права читать список пользователей.",
		"health_users_unavailable":          "Не удалось загрузить список пользователей.",
		"health_sessions_forbidden":         "У XMPP Admin нет права читать подключения.",
		"health_sessions_unavailable":       "Не удалось загрузить активные подключения.",
		"health_rooms_forbidden":            "У XMPP Admin нет права читать список комнат.",
		"health_rooms_unavailable":          "Не удалось загрузить активные комнаты.",
		"users":                             "Пользователи",
		"sessions":                          "Сессии",
		"rooms":                             "Комнаты",
		"health":                            "Состояние",
		"users_subtitle":                    "Найдите аккаунт и проверьте его подключение к %s",
		"sessions_subtitle":                 "Активные подключения устройств к %s",
		"rooms_subtitle":                    "Активные групповые чаты на сервере %s",
		"health_subtitle":                   "Что работает и что требует внимания на %s",
		"data_unavailable":                  "Не удалось получить данные из ejabberd",
		"sessions_lookup_unavailable":       "Состояние сессий временно недоступно; список зарегистрированных аккаунтов всё равно показан.",
		"registered_accounts":               "Аккаунты",
		"account":                           "Аккаунт",
		"jid":                               "JID",
		"connection":                        "Подключение",
		"connected":                         "подключён",
		"not_connected":                     "не подключён",
		"online":                            "в сети",
		"offline":                           "не в сети",
		"no_users":                          "Зарегистрированных аккаунтов нет",
		"active_sessions":                   "Активные подключения",
		"resource":                          "Устройство / ресурс",
		"full_jid":                          "Полный JID",
		"no_sessions":                       "Активных сессий нет",
		"online_rooms":                      "Активные комнаты",
		"room":                              "Комната",
		"service":                           "Сервис",
		"no_rooms":                          "Онлайн-комнат нет",
		"infrastructure_health":             "Состояние сервера",
		"registered_users_metric":           "Аккаунтов",
		"online_sessions_metric":            "Подключений",
		"online_rooms_metric":               "Активных комнат (сервер)",
		"api_latency":                       "Задержка API",
		"ejabberd_api_latency":              "Задержка API ejabberd",
		"checks":                            "Сервисы и проверки",
		"component":                         "Компонент",
		"detail":                            "Подробности",
		"healthy":                           "работает",
		"unavailable":                       "недоступно",
		"xmpp_admin_service":                "XMPP Admin",
		"xmpp_admin_running":                "Панель отвечает.",
		"configuration":                     "Настройки сервера",
		"read_from_config":                  "Настройки загружены.",
		"ejabberd_api":                      "Связь с ejabberd",
		"status_command_ok":                 "Связь с ejabberd работает.",
		"read_access_ok":                    "Данные доступны.",
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

func humanRaw(raw any) string {
	return strings.TrimSpace(fmt.Sprint(raw))
}

func humanAccessRule(lang string, raw any) string {
	value := strings.ToLower(humanRaw(raw))
	switch value {
	case "all", "allow", "true":
		return tr(lang, "all_users")
	case "none", "deny", "false", "":
		return tr(lang, "no_one")
	default:
		return tr(lang, "server_policy")
	}
}

func humanLimit(lang string, raw any) string {
	value := strings.ToLower(humanRaw(raw))
	switch value {
	case "infinity", "infinite", "unlimited":
		return tr(lang, "no_limit")
	case "":
		return "—"
	default:
		return humanRaw(raw)
	}
}

func humanAuto(lang string, raw any) string {
	value := strings.ToLower(humanRaw(raw))
	switch value {
	case "", "auto", "automatic":
		return tr(lang, "automatic")
	case "none", "disabled", "off":
		return tr(lang, "not_used")
	default:
		return tr(lang, "configured")
	}
}

func humanDuration(lang string, raw any) string {
	value := humanRaw(raw)
	seconds, err := strconv.ParseInt(value, 10, 64)
	if err != nil || seconds < 0 {
		if value == "" {
			return "—"
		}
		return value
	}
	if seconds == 0 {
		return "0"
	}
	if seconds%86400 == 0 {
		return humanDurationUnit(lang, seconds/86400, "day")
	}
	if seconds%3600 == 0 {
		return humanDurationUnit(lang, seconds/3600, "hour")
	}
	if seconds%60 == 0 {
		return humanDurationUnit(lang, seconds/60, "minute")
	}
	return humanDurationUnit(lang, seconds, "second")
}

func humanDurationUnit(lang string, value int64, unit string) string {
	if lang == langRU {
		var one, few, many string
		switch unit {
		case "day":
			one, few, many = "день", "дня", "дней"
		case "hour":
			one, few, many = "час", "часа", "часов"
		case "minute":
			one, few, many = "минута", "минуты", "минут"
		default:
			one, few, many = "секунда", "секунды", "секунд"
		}
		mod10, mod100 := value%10, value%100
		word := many
		if mod10 == 1 && mod100 != 11 {
			word = one
		} else if mod10 >= 2 && mod10 <= 4 && (mod100 < 12 || mod100 > 14) {
			word = few
		}
		return fmt.Sprintf("%d %s", value, word)
	}
	label := unit
	if value != 1 {
		label += "s"
	}
	return fmt.Sprintf("%d %s", value, label)
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
