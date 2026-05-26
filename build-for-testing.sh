#!/bin/bash

# Ensure script exits on error
set -e

echo "=== 1. Building local Goto binary ==="
go build -ldflags="-X 'goto/src/cmd.VersionGoto=1.1.0-testing'" -o ./goto src/main.go
echo "Goto binary compiled successfully."

echo ""
echo "=== 2. Running uninstall.sh to clean up previous installations ==="
if [ -f "./uninstall.sh" ]; then
    # Run the uninstallation script
    # We can pass 'y' automatically to remove configuration or let the user choose
    ./uninstall.sh
else
    echo "Warning: uninstall.sh not found, skipping uninstallation step."
fi

echo ""
echo "=== 3. Running local installation (replicating install.sh) ==="
INSTALL_DIR="${HOME}/.local/bin"
mkdir -p "${INSTALL_DIR}"

BINARY_PATH="${INSTALL_DIR}/goto"
echo "Installing Goto binary to ${BINARY_PATH}..."
cp ./goto "${BINARY_PATH}"
chmod +x "${BINARY_PATH}"

# Detect shell
if [ -n "$SHELL" ]; then
    DETECTED_SHELL=$(basename "$SHELL")
else
    DETECTED_SHELL="profile (default)"
fi
echo "Detected shell: $DETECTED_SHELL"

# Run goto init
echo "Running 'goto init'..."
"${BINARY_PATH}" init