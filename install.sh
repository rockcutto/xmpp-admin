#!/usr/bin/env bash
set -euo pipefail

SERVICE_NAME="xmpp-admin"
INSTALL_BIN="/usr/local/bin/xmpp-admin"
CONFIG_DIR="/etc/xmpp-admin"
ENV_FILE="${CONFIG_DIR}/xmpp-admin.env"
UNIT_FILE="/etc/systemd/system/${SERVICE_NAME}.service"
DEFAULT_LISTEN="127.0.0.1:8090"

log() { printf '[xmpp-admin] %s\n' "$*"; }
die() { printf '[xmpp-admin] ERROR: %s\n' "$*" >&2; exit 1; }

if [ "$(id -u)" -ne 0 ]; then
  die "run as root: sudo bash install.sh"
fi

SCRIPT_DIR="$(cd -- "$(dirname -- "$0")" >/dev/null 2>&1 && pwd)"
ADMIN_DIR="${SCRIPT_DIR}"

if [ ! -f "${ADMIN_DIR}/go.mod" ]; then
  die "admin/go.mod not found"
fi

if [ -r /etc/os-release ]; then
  # shellcheck source=/dev/null
  . /etc/os-release
  case "${ID:-}" in
    ubuntu|debian) ;;
    *) log "warning: this installer is tested for Ubuntu/Debian, detected ${ID:-unknown}" ;;
  esac
fi

if ! command -v go >/dev/null 2>&1; then
  log "installing Go toolchain"
  export DEBIAN_FRONTEND=noninteractive
  apt-get update
  apt-get install -y --no-install-recommends golang-go ca-certificates
fi

if ! command -v curl >/dev/null 2>&1; then
  export DEBIAN_FRONTEND=noninteractive
  apt-get update
  apt-get install -y --no-install-recommends curl ca-certificates
fi

GO_VERSION="$(go env GOVERSION 2>/dev/null || true)"
GO_VERSION_NUM="${GO_VERSION#go}"
log "building with ${GO_VERSION:-go}"
if ! dpkg --compare-versions "${GO_VERSION_NUM:-0}" ge "1.22"; then
  die "Go 1.22 or newer is required; found ${GO_VERSION:-unknown}"
fi

(
  cd "${ADMIN_DIR}"
  log "running tests"
  GOFLAGS="-mod=readonly" go test ./...
  log "running go vet"
  GOFLAGS="-mod=readonly" go vet ./...
  log "building binary"
  CGO_ENABLED=0 GOFLAGS="-mod=readonly" go build -buildvcs=false -trimpath -ldflags="-s -w" -o /tmp/xmpp-admin .
)
install -o root -g root -m 0755 /tmp/xmpp-admin "${INSTALL_BIN}"
rm -f /tmp/xmpp-admin

if ! id xmpp-admin >/dev/null 2>&1; then
  log "creating system user xmpp-admin"
  useradd --system --create-home --home-dir /var/lib/xmpp-admin --shell /usr/sbin/nologin xmpp-admin
fi
if getent group ejabberd >/dev/null 2>&1; then
  usermod -a -G ejabberd xmpp-admin
fi

detect_ejabberd_config() {
  if [ -n "${EJABBERD_CONFIG:-}" ] && [ -r "${EJABBERD_CONFIG}" ]; then
    printf '%s' "${EJABBERD_CONFIG}"
    return
  fi
  for p in /etc/ejabberd/ejabberd.yml /opt/ejabberd/conf/ejabberd.yml; do
    if [ -r "$p" ]; then
      printf '%s' "$p"
      return
    fi
  done
  printf '%s' "/etc/ejabberd/ejabberd.yml"
}

EJABBERD_CONFIG_PATH="$(detect_ejabberd_config)"
log "ejabberd config: ${EJABBERD_CONFIG_PATH}"

mkdir -p "${CONFIG_DIR}"
chmod 0750 "${CONFIG_DIR}"

random_secret() {
  if command -v openssl >/dev/null 2>&1; then
    openssl rand -base64 36 | tr -d '\n'
  else
    head -c 48 /dev/urandom | base64 | tr -d '\n'
  fi
}

existing_value() {
  key="$1"
  if [ -f "${ENV_FILE}" ]; then
    sed -n "s/^${key}=//p" "${ENV_FILE}" | tail -n 1
  fi
}

ADMIN_USER_VALUE="${ADMIN_USER:-$(existing_value ADMIN_USER)}"
ADMIN_USER_VALUE="${ADMIN_USER_VALUE:-admin}"

ADMIN_PASSWORD_VALUE="${ADMIN_PASSWORD:-$(existing_value ADMIN_PASSWORD)}"
GENERATED_ADMIN_PASSWORD=0
if [ -z "${ADMIN_PASSWORD_VALUE}" ]; then
  ADMIN_PASSWORD_VALUE="$(random_secret)"
  GENERATED_ADMIN_PASSWORD=1
fi

API_USER_VALUE="${EJABBERD_API_USER:-$(existing_value EJABBERD_API_USER)}"
API_PASSWORD_VALUE="${EJABBERD_API_PASSWORD:-$(existing_value EJABBERD_API_PASSWORD)}"

if [ -t 0 ] && [ -z "${API_USER_VALUE}" ]; then
  printf 'ejabberd API JID (leave empty to configure later): '
  read -r API_USER_VALUE
fi
if [ -t 0 ] && [ -n "${API_USER_VALUE}" ] && [ -z "${API_PASSWORD_VALUE}" ]; then
  printf 'ejabberd API password: '
  read -r -s API_PASSWORD_VALUE
  printf '\n'
fi

XMPP_DOMAIN_VALUE="${XMPP_DOMAIN:-$(existing_value XMPP_DOMAIN)}"
API_URL_VALUE="${EJABBERD_API:-$(existing_value EJABBERD_API)}"
BASE_PATH_VALUE="${ADMIN_BASE_PATH:-$(existing_value ADMIN_BASE_PATH)}"
BASE_PATH_VALUE="${BASE_PATH_VALUE:-/xmpp-admin}"

umask 077
{
  printf 'ADMIN_LISTEN=%s\n' "${ADMIN_LISTEN:-${DEFAULT_LISTEN}}"
  printf 'ADMIN_BASE_PATH=%s\n' "${BASE_PATH_VALUE}"
  printf 'ADMIN_USER=%s\n' "${ADMIN_USER_VALUE}"
  printf 'ADMIN_PASSWORD=%s\n' "${ADMIN_PASSWORD_VALUE}"
  printf 'EJABBERD_CONFIG=%s\n' "${EJABBERD_CONFIG_PATH}"
  [ -n "${XMPP_DOMAIN_VALUE}" ] && printf 'XMPP_DOMAIN=%s\n' "${XMPP_DOMAIN_VALUE}"
  [ -n "${API_URL_VALUE}" ] && printf 'EJABBERD_API=%s\n' "${API_URL_VALUE}"
  [ -n "${API_USER_VALUE}" ] && printf 'EJABBERD_API_USER=%s\n' "${API_USER_VALUE}"
  [ -n "${API_PASSWORD_VALUE}" ] && printf 'EJABBERD_API_PASSWORD=%s\n' "${API_PASSWORD_VALUE}"
} > "${ENV_FILE}"
chown root:xmpp-admin "${ENV_FILE}"
chmod 0640 "${ENV_FILE}"

cat > "${UNIT_FILE}" <<'UNIT'
[Unit]
Description=XMPP Admin
After=network-online.target ejabberd.service
Wants=network-online.target

[Service]
Type=simple
User=xmpp-admin
Group=xmpp-admin
EnvironmentFile=/etc/xmpp-admin/xmpp-admin.env
ExecStart=/usr/local/bin/xmpp-admin
Restart=on-failure
RestartSec=3
TimeoutStopSec=10
UMask=0077

NoNewPrivileges=true
PrivateTmp=true
PrivateDevices=true
ProtectHome=true
ProtectSystem=strict
ProtectControlGroups=true
ProtectKernelModules=true
ProtectKernelTunables=true
ProtectKernelLogs=true
RestrictSUIDSGID=true
LockPersonality=true
MemoryDenyWriteExecute=true
CapabilityBoundingSet=
AmbientCapabilities=
RestrictAddressFamilies=AF_UNIX AF_INET AF_INET6
SystemCallArchitectures=native
ReadOnlyPaths=-/etc/ejabberd -/opt/ejabberd/conf

[Install]
WantedBy=multi-user.target
UNIT

detect_https_server_name() {
  if ! command -v nginx >/dev/null 2>&1; then
    return 0
  fi
  nginx -T 2>/dev/null \
    | awk '
      /listen[[:space:]].*443/ { https=1 }
      https && /server_name[[:space:]]/ {
        for (i=2; i<=NF; i++) {
          gsub(/;/, "", $i)
          if ($i != "_" && $i !~ /^\$/) { print $i; exit }
        }
      }
      /}/ { if (https) https=0 }
    '
}

systemctl daemon-reload
systemctl enable "${SERVICE_NAME}.service" >/dev/null
systemctl restart "${SERVICE_NAME}.service"

sleep 1
log "service status:"
systemctl --no-pager --full status "${SERVICE_NAME}.service" | sed -n '1,12p' || true
if ! systemctl is-active --quiet "${SERVICE_NAME}.service"; then
  journalctl -u "${SERVICE_NAME}.service" -n 40 --no-pager || true
  die "xmpp-admin service did not start"
fi

if command -v nginx >/dev/null 2>&1; then
  log "nginx detected: $(nginx -v 2>&1)"
  HTTPS_SERVER_NAME="$(detect_https_server_name || true)"
  NGINX_SNIPPET="${CONFIG_DIR}/nginx-xmpp-admin.conf"
  cat > "${NGINX_SNIPPET}" <<'NGINX'
# Include inside an existing HTTPS server {} block.
# XMPP Admin remains bound to 127.0.0.1:8090.
location = /xmpp-admin {
    return 308 /xmpp-admin/;
}

# Readiness performs a live ejabberd API check and is intentionally local-only.
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
NGINX
  chmod 0644 "${NGINX_SNIPPET}"
  log "nginx helper written to ${NGINX_SNIPPET}"
  if [ -n "${HTTPS_SERVER_NAME}" ]; then
    log "detected HTTPS vhost: ${HTTPS_SERVER_NAME}"
  fi
  if nginx -t >/dev/null 2>&1; then
    log "existing nginx configuration: OK"
  else
    log "warning: existing nginx configuration does not pass nginx -t; it was not modified"
  fi
else
  log "nginx not detected; admin remains available only on ${DEFAULT_LISTEN}"
fi

HTTP_CODE="$(curl -sS -o /tmp/xmpp-admin-health.json -w '%{http_code}' "http://127.0.0.1:8090${BASE_PATH_VALUE}/healthz" || true)"
log "liveness HTTP ${HTTP_CODE}"
if [ -s /tmp/xmpp-admin-health.json ]; then
  cat /tmp/xmpp-admin-health.json
  printf '\n'
fi
rm -f /tmp/xmpp-admin-health.json
if [ "${HTTP_CODE}" != "200" ]; then
  die "xmpp-admin liveness check failed"
fi

READY_CODE="$(curl -sS -o /tmp/xmpp-admin-ready.json -w '%{http_code}' "http://127.0.0.1:8090${BASE_PATH_VALUE}/readyz" || true)"
log "ejabberd readiness HTTP ${READY_CODE}"
if [ "${READY_CODE}" != "200" ]; then
  log "warning: XMPP Admin is running, but ejabberd invite integration is not ready yet"
fi
if [ -s /tmp/xmpp-admin-ready.json ]; then
  cat /tmp/xmpp-admin-ready.json
  printf '\n'
fi
rm -f /tmp/xmpp-admin-ready.json

log "browser login uses ADMIN_USER/ADMIN_PASSWORD; EJABBERD_API_* are backend-only credentials"

if [ "${GENERATED_ADMIN_PASSWORD}" -eq 1 ]; then
  printf '\n'
  log "generated admin credentials (save them now):"
  printf '  user: %s\n' "${ADMIN_USER_VALUE}"
  printf '  password: %s\n' "${ADMIN_PASSWORD_VALUE}"
fi

if [ -z "${API_USER_VALUE}" ] || [ -z "${API_PASSWORD_VALUE}" ]; then
  printf '\n'
  log "ejabberd API credentials are not configured yet."
  log "set EJABBERD_API_USER and EJABBERD_API_PASSWORD in ${ENV_FILE}, then:"
  printf '  sudo systemctl restart %s\n' "${SERVICE_NAME}"
fi

printf '\n'
log "installed: ${INSTALL_BIN}"
log "environment: ${ENV_FILE}"
log "logs: journalctl -u ${SERVICE_NAME} -f"
log "local UI: http://${DEFAULT_LISTEN}${BASE_PATH_VALUE}/admin"
log "installation complete: service is live, config discovered, nginx checked"
if command -v nginx >/dev/null 2>&1; then
  log "nginx was not modified; helper: ${CONFIG_DIR}/nginx-xmpp-admin.conf"
fi
