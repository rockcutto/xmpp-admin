# XMPP Admin

A small, neutral, self-hosted administration panel for **ejabberd**.

The current release focuses on native account invitations provided by `mod_invites`. XMPP Admin does not create a second invitation database: ejabberd remains the authority for invite tokens, validity, limits, reserved usernames, registration and inviter/invitee state.

The UI is product-neutral, has English and Russian locales, and uses a self-contained Material 3-inspired light/dark interface.

> Project status: pre-release. The invitation workflow is implemented; users, sessions, rooms and infrastructure health are planned next.

## What it does

- Lists native ejabberd invitations.
- Creates native account invitations with an optional preselected username.
- Revokes a single invite by expiring its token.
- Shows invites created by normal registered users in XMPP clients as well as invites created from the panel.
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

XMPP Admin needs only these ejabberd API commands:

```text
list_invites
generate_invite
generate_invite_with_username
expire_invite_by_token
```

Do **not** grant `"*"`.

See [SECURITY.md](SECURITY.md) before publishing an Internet-facing installation.

## Requirements

Recommended:

- Ubuntu 24.04+ with systemd, or another Debian-family release whose packaged Go is 1.22+;
- an existing ejabberd installation with `mod_invites`;
- nginx or another HTTPS reverse proxy;
- Go 1.22+ when building from source.

The installer uses `apt`/`dpkg`; on an older Debian-family release, install Go 1.22+ yourself before running it.

The integration uses the native `mod_invites` command API introduced in ejabberd 26.01. Development was validated against the 26.07 API shape; for production, use a currently supported and patched ejabberd release. Live ejabberd integration is still part of the release checklist rather than CI.

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
5. creates the `xmpp-admin` system user;
6. finds the existing ejabberd config;
7. creates `/etc/xmpp-admin/xmpp-admin.env`;
8. installs and starts a hardened systemd service;
9. keeps the web service on `127.0.0.1:8090`;
10. detects nginx and writes a reverse-proxy example without modifying nginx automatically.

The installer is safe to rerun and does **not** rewrite `ejabberd.yml` or restart ejabberd.

### Generated admin password

On the first installation, if `ADMIN_PASSWORD` is not supplied, the installer generates a random password and prints it once. Save it.

The environment file is stored at:

```text
/etc/xmpp-admin/xmpp-admin.env
```

with restrictive permissions.

## 1. Configure the ejabberd API

If your existing ejabberd config already has a suitable local `mod_http_api` listener and dedicated API identity, XMPP Admin can use it.

Otherwise merge the relevant parts from:

```text
deploy/ejabberd.example.yml
```

into your existing `ejabberd.yml`.

Example:

```yaml
listen:
  -
    port: 5281
    module: ejabberd_http
    ip: 127.0.0.1
    request_handlers:
      /api: mod_http_api

acl:
  xmpp_admin_api:
    user:
      - "xmpp-admin@example.org"

api_permissions:
  "XMPP Admin invite API":
    from:
      - mod_http_api
    who:
      acl: xmpp_admin_api
    what:
      - generate_invite
      - generate_invite_with_username
      - list_invites
      - expire_invite_by_token
```

Replace `example.org` with your XMPP domain.

Create a dedicated API account using your normal ejabberd administration procedure. Avoid placing the password directly in shell history:

```bash
read -rsp 'ejabberd API password: ' API_PASS; echo
sudo ejabberdctl register xmpp-admin example.org "$API_PASS"
unset API_PASS
```

Then edit:

```bash
sudoedit /etc/xmpp-admin/xmpp-admin.env
```

and set:

```text
EJABBERD_API_USER=xmpp-admin@example.org
EJABBERD_API_PASSWORD=A-LONG-RANDOM-PASSWORD
```

Restart only XMPP Admin:

```bash
sudo systemctl restart xmpp-admin
```

Validate your ejabberd configuration with the tools appropriate for your installation before reloading it.

## 2. Configure nginx

The installer writes a helper to:

```text
/etc/xmpp-admin/nginx-xmpp-admin.conf
```

A repository copy is also available at:

```text
deploy/nginx.example.conf
```

The default base path is `/xmpp-admin`.

Add the locations to an existing HTTPS `server {}` block:

```nginx
location = /xmpp-admin {
    return 308 /xmpp-admin/;
}

location /xmpp-admin/ {
    client_max_body_size 64k;

    proxy_pass http://127.0.0.1:8090;
    proxy_http_version 1.1;

    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_set_header X-Forwarded-Host $host;

    proxy_read_timeout 30s;
}
```

Check and reload nginx:

```bash
sudo nginx -t
sudo systemctl reload nginx
```

Then open:

```text
https://your-host.example/xmpp-admin/admin
```

No new public port is required.

## 3. Verify the installation

Service status:

```bash
sudo systemctl status xmpp-admin
```

Logs:

```bash
sudo journalctl -u xmpp-admin -f
```

Local health endpoint:

```bash
curl http://127.0.0.1:8090/xmpp-admin/healthz
```

Expected response:

```json
{"ok":true}
```

The health endpoint is liveness-only and intentionally does not disclose the XMPP domain, config path or API state.

A second endpoint, `/readyz`, performs a live `list_invites` check against ejabberd. It is intended for local installation diagnostics only and the supplied nginx configuration does not expose it.

## Updating

```bash
cd xmpp-admin
git pull --ff-only
sudo bash install.sh
```

The installer preserves existing credentials from `/etc/xmpp-admin/xmpp-admin.env`.

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
- New write actions must preserve same-origin/CSRF protection.
- Keep the ejabberd command whitelist narrow.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## Vulnerability reports

See [SECURITY.md](SECURITY.md). Please do not disclose security vulnerabilities in public issues before a fix is available.

## License

MIT. See [LICENSE](LICENSE).
