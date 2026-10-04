package main

const opsTemplate = `
{{define "ops"}}
<!doctype html>
<html lang="{{.Lang}}">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>XMPP Admin — {{tr .Lang .TitleKey}}</title>
<link href="{{p "/assets/material3.css"}}" rel="stylesheet">
<script src="{{p "/assets/theme.js"}}"></script>
</head>
<body>
<div class="page">
<aside class="nav-rail d-print-none" aria-label="{{tr .Lang "admin_sections"}}">
  <a class="nav-rail-brand" href="{{p "/admin"}}" aria-label="XMPP Admin" title="XMPP Admin">
    <span aria-hidden="true">XA</span>
  </a>
  <nav class="nav-rail-items">
    <a class="nav-rail-link" href="{{p "/admin"}}" aria-label="{{tr .Lang "invitations"}}" title="{{tr .Lang "invitations"}}">
      <span class="nav-rail-icon" aria-hidden="true"><svg viewBox="0 0 24 24"><path d="M4 5h16v14H4z"/><path d="m4 7 8 6 8-6"/></svg></span>
      <span class="nav-rail-label">{{tr .Lang "nav_invites"}}</span>
    </a>
    <a class="nav-rail-link {{if eq .Active "users"}}active{{end}}" href="{{p "/admin/users"}}" aria-label="{{tr .Lang "users"}}" title="{{tr .Lang "users"}}">
      <span class="nav-rail-icon" aria-hidden="true"><svg viewBox="0 0 24 24"><path d="M16 20v-2a4 4 0 0 0-4-4H7a4 4 0 0 0-4 4v2"/><circle cx="9.5" cy="7" r="4"/><path d="M17 11a4 4 0 0 1 4 4v5"/><path d="M16 3.3a4 4 0 0 1 0 7.4"/></svg></span>
      <span class="nav-rail-label">{{tr .Lang "nav_users"}}</span>
    </a>
    <a class="nav-rail-link {{if eq .Active "sessions"}}active{{end}}" href="{{p "/admin/sessions"}}" aria-label="{{tr .Lang "sessions"}}" title="{{tr .Lang "sessions"}}">
      <span class="nav-rail-icon" aria-hidden="true"><svg viewBox="0 0 24 24"><circle cx="12" cy="12" r="8"/><path d="M12 8v5l3 2"/></svg></span>
      <span class="nav-rail-label">{{tr .Lang "nav_sessions"}}</span>
    </a>
    <a class="nav-rail-link {{if eq .Active "rooms"}}active{{end}}" href="{{p "/admin/rooms"}}" aria-label="{{tr .Lang "rooms"}}" title="{{tr .Lang "rooms"}}">
      <span class="nav-rail-icon" aria-hidden="true"><svg viewBox="0 0 24 24"><path d="M4 5h12a3 3 0 0 1 3 3v5a3 3 0 0 1-3 3H9l-5 4z"/><path d="M8 9h7M8 12h5"/></svg></span>
      <span class="nav-rail-label">{{tr .Lang "nav_rooms"}}</span>
    </a>
    <a class="nav-rail-link {{if eq .Active "health"}}active{{end}}" href="{{p "/admin/health"}}" aria-label="{{tr .Lang "health"}}" title="{{tr .Lang "health"}}">
      <span class="nav-rail-icon" aria-hidden="true"><svg viewBox="0 0 24 24"><path d="M3 12h4l2-5 4 10 2-5h6"/><path d="M5 4h14a2 2 0 0 1 2 2v12a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2z"/></svg></span>
      <span class="nav-rail-label">{{tr .Lang "nav_health"}}</span>
    </a>
  </nav>
</aside>
<div class="app-main">
<header class="navbar navbar-expand-md d-print-none">
  <div class="container-xl">
    <a class="navbar-brand navbar-brand-autodark" href="{{p "/admin"}}">XMPP Admin</a>
    <div class="navbar-nav flex-row order-md-last m3-toolbar-actions">
      <button class="btn m3-theme-toggle" type="button" data-theme-toggle data-label-dark="{{tr .Lang "theme_dark"}}" data-label-light="{{tr .Lang "theme_light"}}" aria-label="{{tr .Lang "theme_dark"}}" title="{{tr .Lang "theme_dark"}}">
        <span class="m3-theme-icon" aria-hidden="true"></span>
      </button>
      <div class="btn-group btn-group-sm" role="group" aria-label="Language">
        <a class="btn {{if eq .Lang "en"}}btn-primary{{else}}btn-outline-secondary{{end}}" href="{{p "/lang"}}?lang=en&next={{.CurrentPath}}">EN</a>
        <a class="btn {{if eq .Lang "ru"}}btn-primary{{else}}btn-outline-secondary{{end}}" href="{{p "/lang"}}?lang=ru&next={{.CurrentPath}}">RU</a>
      </div>
    </div>
  </div>
</header>



<div class="page-wrapper">
<div class="page-header d-print-none">
  <div class="container-xl">
    <h2 class="page-title">{{tr .Lang .TitleKey}}</h2>
    <div class="text-secondary page-subtitle">{{printf (tr .Lang .SubtitleKey) .Domain}}</div>
  </div>
</div>

<div class="page-body">
<div class="container-xl">

{{if .APIError}}
<div class="alert alert-danger">
  <div class="fw-bold">{{tr .Lang "data_unavailable"}}</div>
  <div class="small mt-1 ops-break">{{.APIError}}</div>
</div>
{{end}}

{{range .Warnings}}
<div class="alert alert-info">{{.}}</div>
{{end}}

{{if eq .Active "users"}}
{{if not .APIError}}
<div class="card">
  <div class="card-header ops-card-header">
    <h3 class="card-title">{{tr .Lang "registered_accounts"}}</h3>
    <span class="badge bg-secondary-lt">{{len .Users}}</span>
  </div>
  <div class="table-responsive">
    <table class="table table-vcenter card-table ops-table responsive-data-table">
      <thead><tr>
        <th>{{tr .Lang "account"}}</th>
        <th>{{tr .Lang "connection"}}</th>
        <th class="text-end">{{tr .Lang "sessions"}}</th>
      </tr></thead>
      <tbody>
      {{range .Users}}
        <tr>
          <td data-label="{{tr $.Lang "account"}}">
            <div class="cell-primary">{{.Username}}</div>
            <div class="cell-secondary">{{.JID}}</div>
          </td>
          <td data-label="{{tr $.Lang "connection"}}">
            {{if eq .SessionState "online"}}<span class="badge bg-green-lt">{{tr $.Lang "connected"}}</span>
            {{else if eq .SessionState "offline"}}<span class="badge bg-secondary-lt">{{tr $.Lang "not_connected"}}</span>
            {{else}}<span class="badge bg-secondary-lt">—</span>{{end}}
          </td>
          <td data-label="{{tr $.Lang "sessions"}}" class="text-end numeric-cell">{{if eq .SessionState "unknown"}}—{{else}}{{.SessionCount}}{{end}}</td>
        </tr>
      {{else}}
        <tr class="empty-row"><td colspan="3" class="text-secondary text-center py-5">{{tr $.Lang "no_users"}}</td></tr>
      {{end}}
      </tbody>
    </table>
  </div>
</div>
{{end}}
{{end}}

{{if eq .Active "sessions"}}
{{if not .APIError}}
<div class="card">
  <div class="card-header ops-card-header">
    <h3 class="card-title">{{tr .Lang "active_sessions"}}</h3>
    <span class="badge bg-green-lt">{{len .Sessions}}</span>
  </div>
  <div class="table-responsive">
    <table class="table table-vcenter card-table ops-table responsive-data-table">
      <thead><tr>
        <th>{{tr .Lang "account"}}</th>
        <th>{{tr .Lang "resource"}}</th>
      </tr></thead>
      <tbody>
      {{range .Sessions}}
        <tr>
          <td data-label="{{tr $.Lang "account"}}">
            <div class="cell-primary">{{.Username}}</div>
            <div class="cell-secondary ops-break">{{.JID}}</div>
          </td>
          <td data-label="{{tr $.Lang "resource"}}">{{if .Resource}}<code class="ops-break">{{.Resource}}</code>{{else}}—{{end}}</td>
        </tr>
      {{else}}
        <tr class="empty-row"><td colspan="2" class="text-secondary text-center py-5">{{tr $.Lang "no_sessions"}}</td></tr>
      {{end}}
      </tbody>
    </table>
  </div>
</div>
{{end}}
{{end}}

{{if eq .Active "rooms"}}
{{if not .APIError}}
<div class="card">
  <div class="card-header ops-card-header">
    <h3 class="card-title">{{tr .Lang "online_rooms"}}</h3>
    <span class="badge bg-green-lt">{{len .Rooms}}</span>
  </div>
  <div class="table-responsive">
    <table class="table table-vcenter card-table ops-table responsive-data-table">
      <thead><tr>
        <th>{{tr .Lang "room"}}</th>
        <th>{{tr .Lang "service"}}</th>
      </tr></thead>
      <tbody>
      {{range .Rooms}}
        <tr>
          <td data-label="{{tr $.Lang "room"}}">
            <div class="cell-primary">{{.Name}}</div>
            <div class="cell-secondary ops-break">{{.JID}}</div>
          </td>
          <td data-label="{{tr $.Lang "service"}}"><code class="ops-break">{{.Service}}</code></td>
        </tr>
      {{else}}
        <tr class="empty-row"><td colspan="2" class="text-secondary text-center py-5">{{tr $.Lang "no_rooms"}}</td></tr>
      {{end}}
      </tbody>
    </table>
  </div>
</div>
{{end}}
{{end}}

{{if eq .Active "health"}}
<div class="metric-grid">
  <div class="metric-card">
    <div class="metric-label">{{tr .Lang "registered_users_metric"}}</div>
    <div class="metric-value">{{.RegisteredMetric}}</div>
  </div>
  <div class="metric-card">
    <div class="metric-label">{{tr .Lang "online_sessions_metric"}}</div>
    <div class="metric-value">{{.OnlineSessionsMetric}}</div>
  </div>
  <div class="metric-card">
    <div class="metric-label">{{tr .Lang "online_rooms_metric"}}</div>
    <div class="metric-value">{{.OnlineRoomsMetric}}</div>
  </div>
  <div class="metric-card">
    <div class="metric-label">{{tr .Lang "ejabberd_api_latency"}}</div>
    <div class="metric-value metric-value-small">{{.APILatency}}</div>
  </div>
</div>

<div class="card mt-3">
  <div class="card-header"><h3 class="card-title">{{tr .Lang "checks"}}</h3></div>
  <div class="table-responsive">
    <table class="table table-vcenter card-table ops-table responsive-data-table health-table">
      <thead><tr>
        <th>{{tr .Lang "component"}}</th>
        <th>{{tr .Lang "status"}}</th>
        <th>{{tr .Lang "detail"}}</th>
      </tr></thead>
      <tbody>
      {{range .Health}}
        <tr>
          <td data-label="{{tr $.Lang "component"}}"><div class="cell-primary">{{.Name}}</div></td>
          <td data-label="{{tr $.Lang "status"}}">
            {{if eq .State "ok"}}<span class="badge bg-green-lt">{{tr $.Lang .StateKey}}</span>
            {{else if eq .State "warn"}}<span class="badge bg-warning-lt">{{tr $.Lang .StateKey}}</span>
            {{else}}<span class="badge bg-red-lt">{{tr $.Lang .StateKey}}</span>{{end}}
          </td>
          <td data-label="{{tr $.Lang "detail"}}" class="ops-break text-secondary">{{.Detail}}</td>
        </tr>
      {{end}}
      </tbody>
    </table>
  </div>
</div>
{{end}}

</div>
</div>
</div>
</div>
</div>
</body>
</html>
{{end}}
`
