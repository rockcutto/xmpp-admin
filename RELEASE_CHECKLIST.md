# Release checklist

Use this before publishing a tag or deployment package.

## Code

- [ ] `go test ./...`
- [ ] `go vet ./...`
- [ ] `bash -n install.sh`
- [ ] `shellcheck install.sh`
- [ ] `govulncheck ./...`
- [ ] `docker build -t xmpp-admin:test .`
- [ ] CI is green on the release commit.

## Security

- [ ] No real passwords, API credentials, invite tokens, private JIDs or production config were committed.
- [ ] No `api_permissions: "*"`.
- [ ] No direct administrative `register` dependency was added to the invite UI.
- [ ] Frontend remains self-contained; no unreviewed third-party JS/CSS.
- [ ] Administrative responses remain `Cache-Control: no-store`.
- [ ] Write actions retain same-origin/CSRF checks.
- [ ] `/healthz` does not expose sensitive server details.
- [ ] `/readyz` is not exposed by the public reverse-proxy configuration.
- [ ] Plain HTTP ejabberd API is loopback-only.

## ejabberd compatibility

- [ ] Tested against the intended ejabberd release.
- [ ] `list_invites` loads existing user-created invitations.
- [ ] active/expired status is rendered correctly.
- [ ] `generate_invite` works.
- [ ] `generate_invite_with_username` works.
- [ ] `expire_invite_by_token` revokes exactly one invite.
- [ ] `include_config_file` deployment was tested if the production config uses it.
- [ ] multi-vhost deployment was tested if `XMPP_DOMAIN` is used.

## Ubuntu deployment

- [ ] Clean install with `sudo bash install.sh`.
- [ ] Re-running installer preserves credentials.
- [ ] systemd service starts as `xmpp-admin`.
- [ ] service listens only on the configured address.
- [ ] nginx example passes `nginx -t` after integration.
- [ ] `curl http://127.0.0.1:8090/xmpp-admin/healthz` returns `{"ok":true}`.
- [ ] EN and RU both render.
- [ ] light/dark/system theme switching works.

## Public repository

- [ ] README.md is current and public documentation is English-only.
- [ ] LICENSE is present.
- [ ] SECURITY.md is present.
- [ ] CONTRIBUTING.md is present.
- [ ] Dependabot config is present.
- [ ] Public CI workflow is present.
- [ ] No references to the parent Android repository remain.
- [ ] Initial public import is a clean snapshot with exactly one root commit and no parent-project history.
- [ ] Repository description/topics are set.
