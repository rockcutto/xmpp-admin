package main

const adminTemplate = `
{{define "admin"}}
<!doctype html>
<html lang="{{.Lang}}">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>XMPP Admin — {{tr .Lang "invitations"}}</title>
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
    <a class="nav-rail-link active" href="{{p "/admin"}}" aria-label="{{tr .Lang "invitations"}}" title="{{tr .Lang "invitations"}}">
      <span class="nav-rail-icon" aria-hidden="true"><svg viewBox="0 0 24 24"><path d="M4 5h16v14H4z"/><path d="m4 7 8 6 8-6"/></svg></span>
      <span class="nav-rail-label">{{tr .Lang "nav_invites"}}</span>
    </a>
    <a class="nav-rail-link" href="{{p "/admin/users"}}" aria-label="{{tr .Lang "users"}}" title="{{tr .Lang "users"}}">
      <span class="nav-rail-icon" aria-hidden="true"><svg viewBox="0 0 24 24"><path d="M16 20v-2a4 4 0 0 0-4-4H7a4 4 0 0 0-4 4v2"/><circle cx="9.5" cy="7" r="4"/><path d="M17 11a4 4 0 0 1 4 4v5"/><path d="M16 3.3a4 4 0 0 1 0 7.4"/></svg></span>
      <span class="nav-rail-label">{{tr .Lang "nav_users"}}</span>
    </a>
    <a class="nav-rail-link" href="{{p "/admin/sessions"}}" aria-label="{{tr .Lang "connections_title"}}" title="{{tr .Lang "connections_title"}}">
      <span class="nav-rail-icon" aria-hidden="true"><svg viewBox="0 0 24 24"><circle cx="12" cy="12" r="8"/><path d="M12 8v5l3 2"/></svg></span>
      <span class="nav-rail-label">{{tr .Lang "nav_sessions"}}</span>
    </a>
    <a class="nav-rail-link" href="{{p "/admin/rooms"}}" aria-label="{{tr .Lang "rooms"}}" title="{{tr .Lang "rooms"}}">
      <span class="nav-rail-icon" aria-hidden="true"><svg viewBox="0 0 24 24"><path d="M4 5h12a3 3 0 0 1 3 3v5a3 3 0 0 1-3 3H9l-5 4z"/><path d="M8 9h7M8 12h5"/></svg></span>
      <span class="nav-rail-label">{{tr .Lang "nav_rooms"}}</span>
    </a>
    <a class="nav-rail-link" href="{{p "/admin/health"}}" aria-label="{{tr .Lang "health"}}" title="{{tr .Lang "health"}}">
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
    <h1 class="page-title">{{tr .Lang "invitations"}}</h1>
    <div class="text-secondary page-subtitle">{{printf (tr .Lang "invitations_subtitle") .Domain}}</div>
  </div>
</div>

<div class="page-body">
<div class="container-xl">

{{if .APIError}}
<div class="alert alert-danger task-error">
  <div class="fw-bold">{{tr .Lang "api_error"}}</div>
  <details class="inline-technical">
    <summary>{{tr .Lang "technical_details"}}</summary>
    <code class="ops-break">{{.APIError}}</code>
  </details>
</div>
{{end}}

<div class="invite-workspace">
  <div class="invite-sidebar">
    <section class="card create-invite-card">
      <div class="card-header"><h2 class="card-title">{{tr .Lang "new_invite"}}</h2></div>
      <div class="card-body">
        <form method="post" action="{{p "/admin/invites/create"}}">
          <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
          <div class="mb-3">
            <label class="form-label">{{tr .Lang "username_optional"}}</label>
            <div class="input-group">
              <input class="form-control" name="username" placeholder="alex" autocomplete="off">
              <span class="input-group-text">@{{.Domain}}</span>
            </div>
            <div class="form-hint">{{tr .Lang "username_optional_hint"}}</div>
          </div>
          <button class="btn btn-primary w-100" type="submit" {{if .APIError}}disabled{{end}}>{{tr .Lang "create_invite"}}</button>
        </form>
        <div class="text-secondary small mt-3">{{tr .Lang "policy_from_ejabberd"}}</div>
      </div>
    </section>

    <section class="card invitation-settings-card">
      <div class="card-header"><h2 class="card-title">{{tr .Lang "server_config"}}</h2></div>
      <div class="card-body">
        <div class="config-summary-grid">
          <div class="config-summary-item">
            <div class="config-label">{{tr .Lang "domain"}}</div>
            <div class="config-value"><code>{{.Domain}}</code></div>
          </div>
          <div class="config-summary-item">
            <div class="config-label">{{tr .Lang "invitations"}}</div>
            <div class="config-value">
              {{if .Server.InvitesEnabled}}<span class="badge bg-green-lt">{{tr .Lang "enabled"}}</span>{{else}}<span class="badge bg-red-lt">{{tr .Lang "disabled"}}</span>{{end}}
            </div>
          </div>
          <div class="config-summary-item">
            <div class="config-label">{{tr .Lang "access_rule"}}</div>
            <div class="config-value">{{humanAccessRule .Lang .Server.InviteAccessRule}}</div>
          </div>
          <div class="config-summary-item">
            <div class="config-label">{{tr .Lang "max_invites"}}</div>
            <div class="config-value">{{humanLimit .Lang .Server.InviteMaxInvites}}</div>
          </div>
          <div class="config-summary-item">
            <div class="config-label">{{tr .Lang "token_ttl"}}</div>
            <div class="config-value">{{humanDuration .Lang .Server.InviteTTLSeconds}}</div>
          </div>
          <div class="config-summary-item">
            <div class="config-label">{{tr .Lang "landing_page"}}</div>
            <div class="config-value">{{humanAuto .Lang .Server.InviteLandingPage}}</div>
          </div>
        </div>

        <details class="technical-details">
          <summary>{{tr .Lang "technical_details"}}</summary>
          <dl class="technical-grid">
            <dt>mod_invites</dt>
            <dd>{{if .Server.InvitesEnabled}}{{tr .Lang "enabled"}}{{else}}{{tr .Lang "disabled"}}{{end}}</dd>
            <dt>{{tr .Lang "access_rule"}}</dt>
            <dd><code>{{if .Server.InviteAccessRule}}{{.Server.InviteAccessRule}}{{else}}none{{end}}</code></dd>
            <dt>{{tr .Lang "max_invites"}}</dt>
            <dd><code>{{.Server.InviteMaxInvites}}</code></dd>
            <dt>{{tr .Lang "token_ttl"}}</dt>
            <dd><code>{{.Server.InviteTTLSeconds}}</code> {{tr .Lang "seconds"}}</dd>
            <dt>{{tr .Lang "templates_dir"}}</dt>
            <dd>{{if .Server.InviteTemplatesDir}}<code>{{.Server.InviteTemplatesDir}}</code>{{else}}{{tr .Lang "ejabberd_default"}}{{end}}</dd>
            <dt>{{tr .Lang "site_name"}}</dt>
            <dd>{{if .Server.InviteSiteName}}{{.Server.InviteSiteName}}{{else}}—{{end}}</dd>
            <dt>{{tr .Lang "db_type"}}</dt>
            <dd><code>{{.Server.InviteDBType}}</code></dd>
            <dt>{{tr .Lang "webchat_url"}}</dt>
            <dd><code>{{.Server.InviteWebchatURL}}</code></dd>
            <dt>mod_register</dt>
            <dd>{{if .Server.RegisterEnabled}}{{tr .Lang "enabled"}}{{else}}{{tr .Lang "disabled"}}{{end}}</dd>
            <dt>{{tr .Lang "allow_modules"}}</dt>
            <dd>{{if .Server.RegisterAllowModules}}<code>{{join .Server.RegisterAllowModules ", "}}</code>{{else}}—{{end}}</dd>
            <dt>{{tr .Lang "config_path"}}</dt>
            <dd><code>{{.ConfigPath}}</code></dd>
          </dl>
        </details>
      </div>
    </section>
  </div>

  <section class="invite-feed-section">
    <div class="list-toolbar">
      {{if .Invites}}
      <div class="segmented-control" role="group" aria-label="{{tr .Lang "status"}}">
        <button class="segment active" type="button" data-invite-filter="all">{{tr .Lang "filter_all"}}</button>
        <button class="segment" type="button" data-invite-filter="active">{{tr .Lang "filter_active"}}</button>
        <button class="segment" type="button" data-invite-filter="expired">{{tr .Lang "filter_expired"}}</button>
      </div>
      {{end}}
      <div class="text-secondary small">{{tr .Lang "native_invites_note"}}</div>
    </div>

    {{if not .APIError}}
    <div class="invite-feed" data-invite-feed>
      {{range $idx, $invite := .Invites}}
      <article class="invite-card {{if and $.Created (eq $idx 0)}}just-created{{end}}" data-invite-state="{{if $invite.Valid}}active{{else}}expired{{end}}">
        {{if and $.Created (eq $idx 0)}}<div class="invite-ready-label">{{tr $.Lang "invite_ready"}}</div>{{end}}
        <div class="invite-card-main">
          <div class="invite-card-heading">
            <div>
              <div class="invite-card-title">{{inviteTypeLabel $.Lang $invite.Type}}</div>
              <div class="invite-card-account">
                {{if $invite.AccountName}}<code>{{$invite.AccountName}}@{{$.Domain}}</code>{{else}}{{tr $.Lang "recipient_chooses_name"}}{{end}}
              </div>
            </div>
            {{if $invite.Valid}}<span class="badge bg-green-lt">{{tr $.Lang "active"}}</span>{{else}}<span class="badge bg-secondary-lt">{{tr $.Lang "expired"}}</span>{{end}}
          </div>

          <div class="invite-meta-grid">
            <div class="invite-meta-item">
              <span class="invite-meta-label">{{tr $.Lang "created_by"}}</span>
              <span class="invite-meta-value ops-break">{{if $invite.Inviter}}{{$invite.Inviter}}{{else}}{{tr $.Lang "server_generated"}}{{end}}</span>
            </div>
            <div class="invite-meta-item">
              <span class="invite-meta-label">{{tr $.Lang "created_when"}}</span>
              <time class="invite-meta-value human-time" datetime="{{$invite.CreatedAt}}" data-human-time="created" title="{{formatTimestamp $invite.CreatedAt}}">{{formatTimestamp $invite.CreatedAt}}</time>
            </div>
            <div class="invite-meta-item">
              <span class="invite-meta-label">{{tr $.Lang "expires_when"}}</span>
              <time class="invite-meta-value human-time" datetime="{{$invite.Expires}}" data-human-time="expiry" title="{{formatTimestamp $invite.Expires}}">{{formatTimestamp $invite.Expires}}</time>
            </div>
          </div>
        </div>

        <div class="invite-card-actions">
          {{with safeURL $invite.LandingPage}}
            <button class="btn btn-sm btn-outline-secondary" type="button" data-copy-value="{{.}}" data-copy-label="{{tr $.Lang "copy_link"}}" data-copied-label="{{tr $.Lang "copied"}}">{{tr $.Lang "copy_link"}}</button>
            <a class="btn btn-sm btn-outline-secondary" href="{{.}}" target="_blank" rel="noreferrer noopener">{{tr $.Lang "open"}}</a>
          {{else}}
            {{with safeURL $invite.TokenURI}}
              <button class="btn btn-sm btn-outline-secondary" type="button" data-copy-value="{{.}}" data-copy-label="{{tr $.Lang "copy_link"}}" data-copied-label="{{tr $.Lang "copied"}}">{{tr $.Lang "copy_link"}}</button>
              <a class="btn btn-sm btn-outline-secondary" href="{{.}}">{{tr $.Lang "open"}}</a>
            {{end}}
          {{end}}
          {{if $invite.Valid}}
          <form method="post" action="{{p "/admin/invites/revoke"}}">
            <input type="hidden" name="csrf_token" value="{{$.CSRFToken}}">
            <input type="hidden" name="token" value="{{$invite.Token}}">
            <button class="btn btn-sm btn-outline-danger" type="submit">{{tr $.Lang "revoke"}}</button>
          </form>
          {{end}}
        </div>
      </article>
      {{else}}
      <div class="empty-state">
        <div class="empty-state-title">{{tr $.Lang "no_invites_title"}}</div>
        <div class="empty-state-body">{{tr $.Lang "no_invites_body"}}</div>
      </div>
      {{end}}
      <div class="empty-state compact filter-empty" data-invite-filter-empty hidden>
        <div class="empty-state-title">{{tr $.Lang "no_filtered_invites"}}</div>
      </div>
    </div>
    {{end}}
  </section>
</div>

</div>
</div>
</div>
</div>
</div>
</body>
</html>
{{end}}
`
