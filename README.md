# XMPP Admin

A small, neutral, self-hosted administration panel for **ejabberd**.

The current pre-release provides native invitation management plus read-only operational views for registered users, active sessions, online MUC rooms and infrastructure health. XMPP Admin does not create a second invitation database: ejabberd remains the authority for invite tokens, validity, limits, reserved usernames, registration and inviter/invitee state.

The UI is product-neutral, has English and Russian locales, and uses a self-contained Material 3-inspired light/dark interface.

> [!WARNING]
> **Pre-release software**
>
> XMPP Admin is under active development and is not yet intended for production-critical environments. APIs, configuration, security assumptions and deployment instructions may still change before the first stable release. Review the configuration carefully and use a dedicated least-privileged ejabberd API account.

## What it does

- Lists native ejabberd invitations.
- Creates native account invitations with an optional preselected username.
- Revokes a single invite by expiring its token.
- Shows invites created by normal registered users in XMPP clients as well as invites created from the panel.
- Lists registered accounts for the selected XMPP vhost.
- Lists active sessions for the selected XMPP vhost.
- Lists server-wide online MUC rooms reported by ejabberd.
- Shows operational health checks and counts without adding write privileges.
- Reads the existing ejabberd configuration instead of duplicating server policy.
- Understands `hosts`, `host_config`, `append_host_config` and `include_config_file` with `allow_only` / `disallow`.
- Displays the effective `mod_invites` and `mod_register` settings for the selected vhost.
- Detects a local `mod_http_api` listener.
- Supports EN/RU and a system/manual day-night theme.
- Contains no external frontend CDN dependencies.

## Security model

Recommended production layout:

```text
Internet
   |
   v
nginx :443 / HTTPS
   |
   v
XMPP Admin 127.0.0.1:8090
   |
   v
ejabberd mod_http_api on loopback/private network
```

The web service and ejabberd API should not be exposed directly to the public Internet.

XMPP Admin needs only this narrow ejabberd API command set:

```text
list_invites
generate_invite
generate_invite_with_username
expire_invite_by_token
registered_users
connected_users
status
muc_online_rooms
```

The first seven commands cover the core panel. `muc_online_rooms` is used by the Rooms page and requires the corresponding MUC administration command to be available; if it is not available, the rest of the panel continues to work and the Rooms/Health pages report it as unavailable.

Do **not** grant `"*"`.

See [SECURITY.md](SECURITY.md) before publishing an Internet-facing installation.

## Requirements

Recommended:

- Ubuntu/Debian-family host with systemd, or an equivalent manual deployment;
- an existing ejabberd installation with `mod_invites`;
- an ejabberd release that exposes the required commands listed below;
- nginx or another HTTPS reverse proxy;
- Go 1.22+ when building from source.

XMPP Admin intentionally does not depend on one exact ejabberd release number. Before deployment, verify that your installation provides the required commands:

```bash
ejabberdctl help list_invites
ejabberdctl help generate_invite
ejabberdctl help generate_invite_with_username
ejabberdctl help expire_invite_by_token
ejabberdctl help registered_users
ejabberdctl help connected_users
ejabberdctl help status
ejabberdctl help muc_online_rooms
```

The `muc_online_rooms` check is optional when you do not use the Rooms page.

If `ejabberdctl` is not in `PATH`, use the control executable shipped with your ejabberd installation. Do not copy a version-specific binary path from another server.

## Quick install on Ubuntu

From an SSH shell:

```bash
git clone https://github.com/rockcutto/xmpp-admin.git
cd xmpp-admin
sudo bash install.sh
```

The installer:

1. installs Go/curl if missing;
2. runs tests and `go vet`;
3. builds the binary;
4. installs it as `/usr/local/bin/xmpp-admin`;
5. creates the unprivileged `xmpp-admin` system user;
6. discovers the existing ejabberd configuration;
7. creates `/etc/xmpp-admin/xmpp-admin.env`;
8. installs and starts a hardened systemd service;
9. keeps the web service on `127.0.0.1:8090`;
10. writes an nginx helper without modifying nginx automatically.

The installer is safe to rerun. It preserves existing credentials in the environment file and does **not** rewrite or restart ejabberd.

### Two separate sets of credentials

XMPP Admin deliberately uses two independent identities:

| Purpose | Variables/account | Used by |
|---|---|---|
| Web panel login | `ADMIN_USER` / `ADMIN_PASSWORD` | Your browser via HTTP Basic Auth |
| ejabberd API service account | `EJABBERD_API_USER` / `EJABBERD_API_PASSWORD` | XMPP Admin backend only |

A normal XMPP account such as an administrator's personal JID is **not** the web-panel login unless you explicitly configure it that way.

On first installation, if `ADMIN_PASSWORD` was not supplied, the installer generates a random web-panel password and prints it once. The persistent values are stored in:

```text
/etc/xmpp-admin/xmpp-admin.env
```

with restrictive permissions.

## 1. Configure the ejabberd API

### Reuse an existing API listener when possible

If ejabberd already has an `ejabberd_http` listener with:

```yaml
request_handlers:
  /api: mod_http_api
```

you do not need to add another listener.

For a TLS listener, use an HTTPS URL whose hostname matches the certificate, for example:

```text
EJABBERD_API=https://xmpp.example.org:5443/api
```

For a loopback-only plaintext listener, use for example:

```text
EJABBERD_API=http://127.0.0.1:5281/api
```

XMPP Admin refuses plaintext HTTP API endpoints on non-loopback addresses. If automatic TLS listener discovery does not match your certificate topology, set `EJABBERD_API` explicitly.

### Create a dedicated API account

Do not reuse a human administrator account. Create a dedicated account such as `xmpp-admin@example.org` and keep its password out of shell history:

```bash
read -rsp 'ejabberd API password: ' API_PASS; echo
sudo ejabberdctl register xmpp-admin example.org "$API_PASS"
unset API_PASS
```

If the account already exists, use your installation's `change_password` command instead.

### Grant only the required commands

Merge the following into the existing ejabberd configuration. Do not replace unrelated ACLs, access rules, listeners, or API permissions.

```yaml
acl:
  xmpp_admin_api:
    user:
      - "xmpp-admin@example.org"

access_rules:
  configure:
    - allow: admin
    - allow: xmpp_admin_api

api_permissions:
  "XMPP Admin API":
    from:
      - mod_http_api
    who:
      acl: xmpp_admin_api
    what:
      - list_invites
      - generate_invite
      - generate_invite_with_username
      - expire_invite_by_token
      - registered_users
      - connected_users
      - status
      - muc_online_rooms
```

The `configure` access rule is required because these API commands carry a `host` argument and ejabberd applies a host-level authorization gate before command execution.

**Keep the two `allow` entries separate.** This is correct:

```yaml
configure:
  - allow: admin
  - allow: xmpp_admin_api
```

This is **not** equivalent:

```yaml
configure:
  allow:
    - admin
    - xmpp_admin_api
```

ACLs grouped inside one `allow` entry are matched together, so the second form can require the caller to match both ACLs and cause an unexpected HTTP 403.

A ready-to-merge example is also provided in:

```text
deploy/ejabberd.example.yml
```

Reload the ejabberd configuration using the control command from your installation:

```bash
sudo ejabberdctl reload_config
```

Then configure the backend credentials:

```bash
sudoedit /etc/xmpp-admin/xmpp-admin.env
```

Set:

```text
EJABBERD_API_USER=xmpp-admin@example.org
EJABBERD_API_PASSWORD=A-LONG-RANDOM-PASSWORD
```

and, when automatic detection is not appropriate:

```text
EJABBERD_API=https://xmpp.example.org:5443/api
```

Restart only XMPP Admin:

```bash
sudo systemctl restart xmpp-admin
```

### Verify backend authorization before exposing the UI

Local liveness:

```bash
curl -i http://127.0.0.1:8090/xmpp-admin/healthz
```

Expected: HTTP 200 with `{"ok":true}`.

Live ejabberd readiness:

```bash
curl -i http://127.0.0.1:8090/xmpp-admin/readyz
```

Expected: HTTP 200 with `{"ok":true}`.

Also verify that the dedicated API account cannot call unrelated **write/destructive** commands. The following test reads the API credentials from the protected environment file and prints only the HTTP status:

```bash
sudo bash -c '
while IFS="=" read -r k v; do
  case "$k" in
    EJABBERD_API) API="$v" ;;
    EJABBERD_API_USER) USER="$v" ;;
    EJABBERD_API_PASSWORD) PASS="$v" ;;
  esac
done < /etc/xmpp-admin/xmpp-admin.env

curl -sS -o /dev/null -w "HTTP %{http_code}\n" \
  --basic --user "$USER:$PASS" \
  -H "Content-Type: application/json" \
  -d "{}" \
  "$API/stop"
'
```

Expected: **HTTP 403**. If `stop` or another command outside the documented allowlist returns 200, stop and tighten the existing `api_permissions` before exposing XMPP Admin.

## 2. Configure nginx

The installer writes a helper to:

```text
/etc/xmpp-admin/nginx-xmpp-admin.conf
```

A repository copy is available at:

```text
deploy/nginx.example.conf
```

Add those `location` blocks **inside the existing HTTPS `server {}` block** for the site:

```nginx
location = /xmpp-admin {
    return 308 /xmpp-admin/;
}

# Local-only readiness probe.
location = /xmpp-admin/readyz {
    return 404;
}

location /xmpp-admin/ {
    client_max_body_size 64k;

    proxy_pass http://127.0.0.1:8090;
    proxy_http_version 1.1;

    proxy_set_header Host $host;
    proxy_set_header Authorization $http_authorization;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_set_header X-Forwarded-Host $host;

    proxy_read_timeout 30s;
}
```

Keep any existing catch-all `location / { ... }` in place; nginx will choose the more specific `/xmpp-admin/` location.

Validate before reload:

```bash
sudo nginx -t
```

Only when the test succeeds:

```bash
sudo systemctl reload nginx
```

Then open:

```text
https://your-host.example/xmpp-admin/admin
```

Use `ADMIN_USER` and `ADMIN_PASSWORD` for the browser login.

## 3. Production smoke test

After the page opens successfully:

1. confirm existing native ejabberd invites are visible;
2. create a normal invite;
3. create an invite with a preselected username;
4. revoke exactly one invite and confirm the other remains;
5. create an invite from a normal XMPP client and confirm it appears after refresh;
6. open Users and confirm registered accounts load;
7. open Sessions and confirm active resources are limited to the selected vhost;
8. open Rooms and confirm online rooms load, or that the page clearly reports `muc_online_rooms` as unavailable;
9. open Health and confirm the API/configuration checks are healthy;
10. re-check that `/xmpp-admin/readyz` is not publicly reachable through nginx.

XMPP Admin reads and writes the native `mod_invites` state; it does not maintain a second invitation database.

## Updating

A source update is not enough by itself. XMPP Admin runs the installed binary at `/usr/local/bin/xmpp-admin`, so after every `git pull` you must rerun the installer to rebuild the binary and restart the systemd service:

```bash
cd xmpp-admin
git pull --ff-only
sudo bash install.sh
```

The installer preserves existing credentials from `/etc/xmpp-admin/xmpp-admin.env`.

After an update that changes forms or CSRF handling, reload the admin page before submitting a form so the browser receives the token generated by the restarted process.

## Configuration

| Variable | Default | Purpose |
|---|---|---|
| `ADMIN_LISTEN` | `127.0.0.1:8090` | HTTP listen address |
| `ADMIN_BASE_PATH` | empty in the binary; installer uses `/xmpp-admin` | Reverse-proxy URL prefix |
| `ADMIN_USER` | `admin` | Basic Auth username |
| `ADMIN_PASSWORD` | required, minimum 12 characters | Basic Auth password |
| `EJABBERD_CONFIG` | auto-detected | Path to the main ejabberd YAML config |
| `XMPP_DOMAIN` | first configured host | Explicit vhost override |
| `EJABBERD_API` | detected from config, else `http://127.0.0.1:5281/api` | ejabberd HTTP API |
| `EJABBERD_API_USER` | empty | Dedicated API account |
| `EJABBERD_API_PASSWORD` | empty | Dedicated API password |

`EJABBERD_API_USER` and `EJABBERD_API_PASSWORD` must either both be set or both be empty.

Plain HTTP for `EJABBERD_API` is accepted only for loopback. Use HTTPS for a non-loopback API endpoint.

See [xmpp-admin.env.example](xmpp-admin.env.example).

## Existing ejabberd configuration

XMPP Admin reads the configuration read-only and currently understands:

- `hosts`;
- `host_config`;
- `append_host_config`;
- recursive `include_config_file`;
- `include_config_file` filters `allow_only` and `disallow`;
- `modules.mod_invites`:
  - `access_create_account`;
  - `max_invites`;
  - `token_expire_seconds`;
  - `landing_page`;
  - `templates_dir`;
  - `site_name`;
  - `db_type`;
  - `webchat_url`;
- `mod_register.allow_modules`;
- a local `ejabberd_http` listener with a `mod_http_api` handler.

The configuration is re-read when the admin page is opened. XMPP Admin does not edit it.

The parser is intentionally a read-only, best-effort view of the documented YAML structures above; ejabberd remains authoritative. If your deployment generates configuration through external preprocessing or unsupported custom constructs, verify the values shown by the panel and use explicit environment overrides where appropriate.

Preselected usernames receive only lightweight safety validation in XMPP Admin. Final XMPP node/JID normalization and acceptance are performed by ejabberd.

For multi-vhost deployments, set `XMPP_DOMAIN` when the first entry in `hosts` is not the host you want to manage.

## Locales and theme

The UI includes:

- English;
- Russian;
- browser `Accept-Language` detection;
- English fallback;
- an EN/RU switch stored in a SameSite cookie;
- light and dark Material 3-inspired themes;
- system theme detection;
- a manual theme switch stored locally in the browser.

All CSS and JavaScript used by the interface are served by XMPP Admin itself.

## Docker

Systemd on the ejabberd host is the recommended deployment when the ejabberd API is loopback-only.

A container build is available:

```bash
docker build -t xmpp-admin .
```

For a Linux host where ejabberd listens only on loopback, host networking is the simplest way to preserve that boundary:

```bash
cp xmpp-admin.env.example .env
chmod 600 .env
# edit .env first

EJABBERD_GID="$(getent group ejabberd | cut -d: -f3)"

docker run --rm \
  --network host \
  --group-add "$EJABBERD_GID" \
  -v /etc/ejabberd:/etc/ejabberd:ro \
  --env-file .env \
  -e ADMIN_LISTEN=127.0.0.1:8090 \
  xmpp-admin
```

This keeps both XMPP Admin and a loopback-only ejabberd API on the host loopback interface. The supplementary ejabberd group is added so the non-root container process can read a typical group-protected configuration directory.

If your configuration uses `include_config_file`, mount the whole configuration directory read-only; mounting only `ejabberd.yml` can hide included files.

For bridge networking or non-Linux Docker hosts, set an explicit `EJABBERD_API=https://...` endpoint reachable through a trusted private network. Do not expose a plaintext ejabberd API publicly just to make container networking easier.

## Development

Run all local checks:

```bash
make check
```

Build:

```bash
make build
```

Individual commands:

```bash
go test ./...
go vet ./...
bash -n install.sh
```

CI also runs:

- Go 1.22 and current Go;
- tests;
- vet;
- shellcheck;
- Docker build;
- Go vulnerability scanning.

## Design boundaries

- ejabberd remains the authority for invite credentials.
- The panel has no invitation database.
- User passwords are never handled or stored by XMPP Admin.
- ejabberd API credentials never reach browser JavaScript.
- Administrative pages are sent with `Cache-Control: no-store`.
- The frontend has no third-party CDN dependency.
- Users, sessions, rooms and health views are read-only.
- New write actions must preserve same-origin/CSRF protection.
- Keep the ejabberd command whitelist narrow.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## Vulnerability reports

See [SECURITY.md](SECURITY.md). Please do not disclose security vulnerabilities in public issues before a fix is available.

## License

MIT. See [LICENSE](LICENSE).
