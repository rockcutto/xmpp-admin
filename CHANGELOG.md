# Changelog

All notable changes to XMPP Admin will be documented in this file.

The project follows semantic versioning once public releases are tagged.

## Unreleased

### Added

- Native ejabberd `mod_invites` listing, generation and single-token revocation.
- Automatic read-only discovery of existing ejabberd configuration.
- Support for vhosts and recursive `include_config_file` configuration.
- English and Russian UI.
- Self-contained Material 3-inspired light/dark interface.
- Ubuntu/Debian installer with hardened systemd service.
- nginx reverse-proxy deployment example.
- Docker build.
- CI across minimum/current Go, shellcheck, container build and Go vulnerability scanning.
- Public security policy, contribution guide and release checklist.
- Separate liveness and local-only ejabberd readiness endpoints.
- Tests for native API wire decoding, config include semantics, CSRF metadata and token-free invite creation redirects.
- Compact invite-table presentation with readable timestamps and human-friendly invite type labels.
- Production deployment guidance based on a live nginx + ejabberd integration.

### Security

- Narrow ejabberd API command allowlist; no wildcard or administrative registration command required by the invite UI.
- Administrative pages use no-store caching, same-origin write checks and restrictive browser security headers.
- Invite URLs are not placed in admin-page query strings after generation.
- Frontend assets are served locally with no third-party CDN runtime dependency.
- Plain HTTP ejabberd API connections are restricted to loopback.
- ejabberd API redirects are disabled and URL-embedded credentials are rejected.
- Native invite tokens are kept out of redirect URLs and browser history.
- Public nginx configuration does not expose the live readiness probe.
- Runtime Basic Auth credential syntax is validated before serving requests.
- Host-scoped invite API access is documented with the required `configure` gate and separate allow entries.
- TLS API autodetection uses the configured XMPP vhost as the certificate hostname instead of a loopback IP.
- nginx examples explicitly forward the Authorization header used by the web-panel Basic Auth boundary.
