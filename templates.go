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
<header class="navbar navbar-expand-md d-print-none">
  <div class="container-xl">
    <a class="navbar-brand navbar-brand-autodark" href="{{p "/admin"}}">XMPP Admin</a>
    <div class="navbar-nav flex-row order-md-last m3-toolbar-actions">
      <button class="btn m3-theme-toggle" type="button" data-theme-toggle data-label-dark="{{tr .Lang "theme_dark"}}" data-label-light="{{tr .Lang "theme_light"}}" aria-label="{{tr .Lang "theme_dark"}}" title="{{tr .Lang "theme_dark"}}">
        <span class="m3-theme-icon" aria-hidden="true"></span>
      </button>
      <div class="btn-group btn-group-sm" role="group" aria-label="Language">
        <a class="btn {{if eq .Lang "en"}}btn-primary{{else}}btn-outline-secondary{{end}}" href="{{p "/lang"}}?lang=en&next=/admin">EN</a>
        <a class="btn {{if eq .Lang "ru"}}btn-primary{{else}}btn-outline-secondary{{end}}" href="{{p "/lang"}}?lang=ru&next=/admin">RU</a>
      </div>
    </div>
  </div>
</header>

<div class="page-wrapper">
<div class="page-header d-print-none">
  <div class="container-xl">
    <h2 class="page-title">{{tr .Lang "invitations"}}</h2>
    <div class="text-secondary">{{printf (tr .Lang "invitations_subtitle") .Domain}}</div>
  </div>
</div>

<div class="page-body">
<div class="container-xl">

{{if .APIError}}
<div class="alert alert-danger">
  <div class="fw-bold">{{tr .Lang "api_error"}}</div>
  <div class="small mt-1">{{.APIError}}</div>
</div>
{{end}}

{{if .Created}}
<div class="alert alert-success">
  <div class="fw-bold mb-1">{{tr .Lang "invite_created"}}</div>
  <div class="small">{{tr .Lang "native_created_note"}}</div>
</div>
{{end}}

<div class="alert alert-info">
  {{tr .Lang "native_invites_note"}}
</div>

<div class="row row-cards">
  <div class="col-lg-4">
    <div class="card">
      <div class="card-header"><h3 class="card-title">{{tr .Lang "new_invite"}}</h3></div>
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
          <button class="btn btn-primary w-100" type="submit" {{if .APIError}}disabled{{end}}>
            {{tr .Lang "create_invite"}}
          </button>
        </form>
        <div class="text-secondary small mt-3">{{tr .Lang "policy_from_ejabberd"}}</div>
      </div>
    </div>

    <div class="card mt-3">
      <div class="card-header"><h3 class="card-title">{{tr .Lang "server_config"}}</h3></div>
      <div class="card-body">
        <dl class="row mb-0">
          <dt class="col-5">{{tr .Lang "domain"}}</dt>
          <dd class="col-7"><code>{{.Domain}}</code></dd>

          <dt class="col-5">mod_invites</dt>
          <dd class="col-7">
            {{if .Server.InvitesEnabled}}
              <span class="badge bg-green-lt">{{tr .Lang "enabled"}}</span>
            {{else}}
              <span class="badge bg-red-lt">{{tr .Lang "disabled"}}</span>
            {{end}}
          </dd>

          <dt class="col-5">{{tr .Lang "access_rule"}}</dt>
          <dd class="col-7"><code>{{if .Server.InviteAccessRule}}{{.Server.InviteAccessRule}}{{else}}none{{end}}</code></dd>

          <dt class="col-5">{{tr .Lang "max_invites"}}</dt>
          <dd class="col-7">{{.Server.InviteMaxInvites}}</dd>

          <dt class="col-5">{{tr .Lang "token_ttl"}}</dt>
          <dd class="col-7">{{.Server.InviteTTLSeconds}} {{tr .Lang "seconds"}}</dd>

          <dt class="col-5">{{tr .Lang "landing_page"}}</dt>
          <dd class="col-7"><code>{{.Server.InviteLandingPage}}</code></dd>

          <dt class="col-5">{{tr .Lang "templates_dir"}}</dt>
          <dd class="col-7">{{if .Server.InviteTemplatesDir}}<code>{{.Server.InviteTemplatesDir}}</code>{{else}}{{tr .Lang "ejabberd_default"}}{{end}}</dd>

          <dt class="col-5">{{tr .Lang "site_name"}}</dt>
          <dd class="col-7">{{if .Server.InviteSiteName}}{{.Server.InviteSiteName}}{{else}}—{{end}}</dd>

          <dt class="col-5">{{tr .Lang "db_type"}}</dt>
          <dd class="col-7"><code>{{.Server.InviteDBType}}</code></dd>

          <dt class="col-5">{{tr .Lang "webchat_url"}}</dt>
          <dd class="col-7"><code>{{.Server.InviteWebchatURL}}</code></dd>

          <dt class="col-5">mod_register</dt>
          <dd class="col-7">
            {{if .Server.RegisterEnabled}}
              <span class="badge bg-green-lt">{{tr .Lang "enabled"}}</span>
            {{else}}
              <span class="badge bg-secondary-lt">{{tr .Lang "disabled"}}</span>
            {{end}}
          </dd>

          <dt class="col-5">{{tr .Lang "allow_modules"}}</dt>
          <dd class="col-7">{{if .Server.RegisterAllowModules}}<code>{{join .Server.RegisterAllowModules ", "}}</code>{{else}}—{{end}}</dd>
        </dl>
        <hr>
        <div class="text-secondary small">
          {{tr .Lang "config_path"}}:<br><code>{{.ConfigPath}}</code>
        </div>
      </div>
    </div>
  </div>

  <div class="col-lg-8">
    <div class="card">
      <div class="table-responsive">
        <table class="table table-vcenter card-table invites-table">
          <thead>
            <tr>
              <th>{{tr .Lang "created"}}</th>
              <th>{{tr .Lang "inviter"}}</th>
              <th>{{tr .Lang "account_name"}}</th>
              <th>{{tr .Lang "expires"}}</th>
              <th>{{tr .Lang "type"}}</th>
              <th>{{tr .Lang "status"}}</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
          {{range .Invites}}
            <tr>
              <td class="text-nowrap invite-date"><time datetime="{{.CreatedAt}}" title="{{.CreatedAt}}">{{formatTimestamp .CreatedAt}}</time></td>
              <td class="invite-jid">{{if .Inviter}}{{.Inviter}}{{else}}{{tr $.Lang "server_generated"}}{{end}}</td>
              <td>{{if .AccountName}}<code>{{.AccountName}}@{{$.Domain}}</code>{{else}}—{{end}}</td>
              <td class="text-nowrap invite-date"><time datetime="{{.Expires}}" title="{{.Expires}}">{{formatTimestamp .Expires}}</time></td>
              <td><span class="badge bg-secondary-lt invite-type">{{inviteTypeLabel $.Lang .Type}}</span></td>
              <td>
                {{if .Valid}}
                  <span class="badge bg-green-lt">{{tr $.Lang "active"}}</span>
                {{else}}
                  <span class="badge bg-secondary-lt">{{tr $.Lang "expired"}}</span>
                {{end}}
              </td>
              <td class="text-end text-nowrap">
                <div class="invite-actions">
                  {{with safeURL .LandingPage}}<a class="btn btn-sm btn-outline-secondary" href="{{.}}" target="_blank" rel="noreferrer noopener">{{tr $.Lang "open"}}</a>{{else}}{{with safeURL .TokenURI}}<a class="btn btn-sm btn-outline-secondary" href="{{.}}">{{tr $.Lang "open"}}</a>{{end}}{{end}}
                  {{if .Valid}}
                  <form method="post" action="{{p "/admin/invites/revoke"}}">
                    <input type="hidden" name="csrf_token" value="{{$.CSRFToken}}">
                    <input type="hidden" name="token" value="{{.Token}}">
                    <button class="btn btn-sm btn-outline-danger" type="submit">{{tr $.Lang "revoke"}}</button>
                  </form>
                  {{end}}
                </div>
              </td>
            </tr>
          {{else}}
            <tr><td colspan="7" class="text-secondary text-center py-5">{{tr $.Lang "no_invites"}}</td></tr>
          {{end}}
          </tbody>
        </table>
      </div>
    </div>
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
