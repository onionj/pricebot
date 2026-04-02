#!/usr/bin/env bash
set -euo pipefail

INSTALL_DIR="/opt/pricebot"
BINARY_NAME="price"
SERVICE_NAME="pricebot"
ENV_FILE="${INSTALL_DIR}/.env"

# --- Helpers ---
info()  { echo -e "\033[1;32m[INFO]\033[0m  $*"; }
error() { echo -e "\033[1;31m[ERROR]\033[0m $*" >&2; }

if [ "$(id -u)" -ne 0 ]; then
    error "This script must be run as root (use sudo)."
    exit 1
fi

# --- Detect architecture ---
ARCH=$(uname -m)
case "$ARCH" in
    x86_64)  GOARCH="amd64" ;;
    aarch64|arm64) GOARCH="arm64" ;;
    *) error "Unsupported architecture: $ARCH"; exit 1 ;;
esac

OS="linux"
ASSET_NAME="${BINARY_NAME}-${OS}-${GOARCH}"

# --- Determine download URL ---
if [ -n "${1:-}" ]; then
    # If a local binary path is provided, use it
    LOCAL_BINARY="$1"
    if [ ! -f "$LOCAL_BINARY" ]; then
        error "File not found: $LOCAL_BINARY"
        exit 1
    fi
    info "Installing from local file: $LOCAL_BINARY"
else
    REPO="onionj/pricebot"
    info "Fetching latest release from GitHub..."
    DOWNLOAD_URL=$(curl -sfL "https://api.github.com/repos/${REPO}/releases/latest" \
        | grep "browser_download_url.*${ASSET_NAME}" \
        | head -1 \
        | cut -d '"' -f 4)

    if [ -z "$DOWNLOAD_URL" ]; then
        error "Could not find release asset: $ASSET_NAME"
        exit 1
    fi

    info "Downloading $DOWNLOAD_URL ..."
    LOCAL_BINARY=$(mktemp)
    curl -sfL -o "$LOCAL_BINARY" "$DOWNLOAD_URL"
fi

# --- Stop existing service if running ---
if systemctl is-active --quiet "$SERVICE_NAME" 2>/dev/null; then
    info "Stopping existing $SERVICE_NAME service..."
    systemctl stop "$SERVICE_NAME"
fi

# --- Install binary ---
mkdir -p "$INSTALL_DIR"
cp "$LOCAL_BINARY" "${INSTALL_DIR}/${BINARY_NAME}"
chmod 755 "${INSTALL_DIR}/${BINARY_NAME}"
info "Binary installed to ${INSTALL_DIR}/${BINARY_NAME}"

# --- Create .env file if it doesn't exist ---
if [ ! -f "$ENV_FILE" ]; then
    cat > "$ENV_FILE" <<'EOF'
BOT_TOKEN=
CHAT_ID=
CHANEL_NAME=
PROXY_LINK=
EOF
    chmod 600 "$ENV_FILE"
    info "Created $ENV_FILE — please edit it with your values."
else
    info "$ENV_FILE already exists, skipping."
fi

# --- Create systemd service ---
cat > "/etc/systemd/system/${SERVICE_NAME}.service" <<EOF
[Unit]
Description=PriceBot Telegram Service
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
WorkingDirectory=${INSTALL_DIR}
ExecStart=${INSTALL_DIR}/${BINARY_NAME}
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable "$SERVICE_NAME"
info "Systemd service created and enabled."

# --- Start or prompt ---
if grep -q '^BOT_TOKEN=$' "$ENV_FILE"; then
    info "Service NOT started — edit $ENV_FILE first, then run:"
    info "  sudo systemctl start $SERVICE_NAME"
else
    systemctl start "$SERVICE_NAME"
    info "Service started."
fi

echo ""
info "Installation complete!"
info "  Binary:  ${INSTALL_DIR}/${BINARY_NAME}"
info "  Config:  ${ENV_FILE}"
info "  Service: sudo systemctl status $SERVICE_NAME"
info "  Logs:    sudo journalctl -u $SERVICE_NAME -f"
