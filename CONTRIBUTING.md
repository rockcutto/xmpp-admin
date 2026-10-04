# Contributing

Contributions are welcome.

## Development

Requirements:

- Go 1.22 or newer;
- Bash;
- optional Docker.

Run the local checks:

```bash
make check
make build
```

Or directly:

```bash
go test ./...
go vet ./...
bash -n install.sh
```

## Design constraints

Please keep these boundaries intact:

- ejabberd remains the authority for users, sessions, rooms and native invite credentials;
- do not introduce a second invite database;
- do not expose ejabberd credentials to browser JavaScript;
- do not broaden `api_permissions` to `"*"`;
- keep the UI client/product neutral;
- keep EN and RU translation catalogs in sync;
- new write actions must have cross-origin/CSRF protection;
- new third-party frontend dependencies require explicit review.

## Pull requests

Keep changes focused and include tests for protocol/config parsing or security-sensitive behavior when practical.
