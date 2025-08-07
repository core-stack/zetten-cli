#!/bin/bash

set -e

CLI_NAME="zetten-cli"
SERVICE_NAME="zetten-service"
INSTALL_DIR="/usr/local/bin"
ROOT_DIR="$HOME/.zetten"
ROOT_CONFIG="$ROOT_DIR/config.yml"

echo "Installing Zetten CLI and Service..."

if [[ ! -f "$CLI_NAME" || ! -f "$SERVICE_NAME" ]]; then
  echo "$CLI_NAME and $SERVICE_NAME not found."
  exit 1
fi

echo "Copying binaries $INSTALL_DIR..."
sudo cp "$CLI_NAME" "$INSTALL_DIR/"
sudo cp "$SERVICE_NAME" "$INSTALL_DIR/"
sudo chmod +x "$INSTALL_DIR/$CLI_NAME"
sudo chmod +x "$INSTALL_DIR/$SERVICE_NAME"

if [[ ! -d "$ROOT_DIR" ]]; then
  echo "Creating configuration directory: $ROOT_DIR"
  mkdir -p "$ROOT_DIR"
fi

if [[ ! -f "$ROOT_CONFIG" ]]; then
  echo "Creating config root default in $ROOT_CONFIG"
  cat > "$ROOT_CONFIG" <<EOF
zettenProjects: []
mirror: []
EOF
fi

echo "Installing service $SERVICE_NAME..."
sudo "$INSTALL_DIR/$SERVICE_NAME" install

echo "Starting service $SERVICE_NAME..."
sudo "$INSTALL_DIR/$SERVICE_NAME" start


echo "Complete installation!"
echo "Configuration in: $ROOT_CONFIG"
echo "Use '$SERVICE_NAME [start|stop|uninstall]' to manage the service"
