#!/bin/bash

# Ensure the script exits on errors
set -e

echo "Starting Goto installation..."

# Determine OS
OS="$(uname -s)"
case "${OS}" in
    Linux*)     OS_NAME=linux;;
    Darwin*)    OS_NAME=darwin;;
    *)          echo "Unsupported OS: ${OS}"; exit 1;;
esac

# Determine Architecture
ARCH="$(uname -m)"
case "${ARCH}" in
    x86_64)     ARCH_NAME=amd64;;
    aarch64)    ARCH_NAME=arm64;;
    arm64)      ARCH_NAME=arm64;;
    *)          echo "Unsupported architecture: ${ARCH}"; exit 1;;
esac

# Create installation directory
INSTALL_DIR="${HOME}/.local/bin"
mkdir -p "${INSTALL_DIR}"

# Get the latest release download URL
echo "Fetching latest release information..."
LATEST_URL=$(curl -s https://api.github.com/repos/Joacohbc/goto/releases/latest | grep "browser_download_url" | grep "goto-${OS_NAME}-${ARCH_NAME}\"" | cut -d '"' -f 4)

if [ -z "${LATEST_URL}" ]; then
    echo "Error: Could not find a binary for ${OS_NAME}-${ARCH_NAME} in the latest release."
    exit 1
fi

# Download the binary
BINARY_PATH="${INSTALL_DIR}/goto"
echo "Downloading Goto from ${LATEST_URL} to ${BINARY_PATH}..."
curl -L -o "${BINARY_PATH}" "${LATEST_URL}"

# Make it executable
chmod +x "${BINARY_PATH}"
echo "Goto successfully downloaded."

# Detect shell
DETECTED_SHELL=$(basename "$SHELL")
if [ -z "$DETECTED_SHELL" ]; then
    DETECTED_SHELL="profile (default)"
fi
echo "Detected shell: $DETECTED_SHELL"

# Run goto init
echo "Running goto init..."
"${BINARY_PATH}" init

echo ""
echo "Installation complete!"
echo "Please restart your terminal or source your shell configuration file to start using 'goto'."
