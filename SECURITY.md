# Security Policy

## Scope

XMPP Admin is an administrative control plane for ejabberd. Security reports are especially important when they involve:

- authentication or authorization bypass;
- CSRF or cross-origin administrative actions;
- exposure of invite tokens or ejabberd API credentials;
- unsafe handling of `ejabberd.yml`;
- command injection or unsafe installer behavior;
- SSRF or unexpected outbound requests;
- privilege escalation in the systemd deployment;
- XSS in the administrative UI.

## Reporting a vulnerability

Please use a GitHub private security advisory for the repository.

Do not open a public issue for an undisclosed vulnerability.

Include, when possible:

- affected version or commit;
- deployment model;
- reproduction steps;
- impact;
- suggested remediation.

Please do not include real account passwords, invite tokens, API credentials, private JIDs, or production configuration secrets in a report.

## Security model

The intended production boundary is:

```text
Internet
   |
   v
nginx / HTTPS
   |
   v
XMPP Admin (loopback)
   |
   v
ejabberd mod_http_api (loopback/private)
```

XMPP Admin should not be exposed directly on a public TCP port. The ejabberd HTTP API should remain loopback-only or use HTTPS on a trusted private network.

The panel intentionally uses a narrow ejabberd command allowlist. Do not grant `"*"`.

`/healthz` is a non-sensitive liveness check. `/readyz` performs a live ejabberd API request and is intended for local diagnostics; the supplied nginx configuration does not publish it.

## Supported deployment

The installer targets Ubuntu/Debian with systemd. Container deployment is supported but requires the operator to provide secure networking between the container and ejabberd.

## Secrets

XMPP Admin does not store user passwords or maintain its own invite database.

Administrative credentials are read from the service environment. Native invite tokens are read from ejabberd only for authenticated administrative operations and pages are returned with `Cache-Control: no-store`.
