#!/bin/bash

# Ensure the script exits on errors
set -e

echo "Starting Goto uninstallation..."

# Remove the binary from ~/.local/bin if it exists
INSTALL_DIR="${HOME}/.local/bin"
BINARY_PATH="${INSTALL_DIR}/goto"

if [ -f "${BINARY_PATH}" ]; then
    echo "Removing binary from ${BINARY_PATH}..."
    rm -f "${BINARY_PATH}"
else
    echo "Binary not found in ${BINARY_PATH}, skipping..."
fi

# Detect shell configuration files
HOME_DIR="${HOME}"
CONFIG_FILES=()

# Common shell configurations
if [ -f "${HOME_DIR}/.bashrc" ]; then CONFIG_FILES+=("${HOME_DIR}/.bashrc"); fi
if [ -f "${HOME_DIR}/.zshrc" ]; then CONFIG_FILES+=("${HOME_DIR}/.zshrc"); fi
if [ -f "${HOME_DIR}/.config/fish/config.fish" ]; then CONFIG_FILES+=("${HOME_DIR}/.config/fish/config.fish"); fi
if [ -f "${HOME_DIR}/.profile" ]; then CONFIG_FILES+=("${HOME_DIR}/.profile"); fi
if [ -f "${HOME_DIR}/.tcshrc" ]; then CONFIG_FILES+=("${HOME_DIR}/.tcshrc"); fi
if [ -f "${HOME_DIR}/.cshrc" ]; then CONFIG_FILES+=("${HOME_DIR}/.cshrc"); fi
if [ -f "${HOME_DIR}/.kshrc" ]; then CONFIG_FILES+=("${HOME_DIR}/.kshrc"); fi

ALIAS_SCRIPT_NAME="alias.sh"

echo "Checking shell configurations for goto alias..."
for FILE in "${CONFIG_FILES[@]}"; do
    if grep -q "${ALIAS_SCRIPT_NAME}" "${FILE}"; then
        echo "Removing goto alias from ${FILE}..."
        # Using sed to remove the source line and the comment
        sed -i.bak '/#Aliases to use goto:/d' "${FILE}"
        sed -i.bak "/source.*${ALIAS_SCRIPT_NAME}/d" "${FILE}"
        rm -f "${FILE}.bak"
    fi
done

# Remove fish completion if it exists
FISH_COMPLETION="${HOME_DIR}/.config/fish/completions/goto.fish"
if [ -f "${FISH_COMPLETION}" ]; then
    echo "Removing fish completion from ${FISH_COMPLETION}..."
    rm -f "${FISH_COMPLETION}"
fi

# Ask user if they want to delete the configuration directory
CONFIG_DIR="${HOME}/.config/goto"
if [ -d "${CONFIG_DIR}" ]; then
    echo ""
    read -p "Do you want to completely remove your goto configuration and saved paths? (${CONFIG_DIR}) [y/N]: " choice
    case "$choice" in
      y|Y )
        echo "Removing configuration directory..."
        rm -rf "${CONFIG_DIR}"
        ;;
      * )
        echo "Keeping configuration directory."
        ;;
    esac
fi

echo ""
echo "Uninstallation complete!"
echo "Please restart your terminal or open a new session to fully apply the changes."
