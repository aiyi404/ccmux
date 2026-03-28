#!/usr/bin/env bash
# ccc installer — builds Go binary and installs to ~/.local/bin
set -euo pipefail

GREEN='\033[32m' YELLOW='\033[33m' RED='\033[31m' RESET='\033[0m'
ok()   { printf "${GREEN}✓${RESET} %s\n" "$*"; }
warn() { printf "${YELLOW}warning:${RESET} %s\n" "$*"; }
die()  { printf "${RED}error:${RESET} %s\n" "$*" >&2; exit 1; }

INSTALL_DIR="${CCC_INSTALL_DIR:-$HOME/.local/bin}"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
BINARY_NAME="ccc"

# check Go
command -v go &>/dev/null || die "Go is required. Install from https://go.dev/dl/"

echo "Building ccmux..."
cd "$SCRIPT_DIR"
go build -ldflags="-s -w" -o "$BINARY_NAME" .
ok "built $BINARY_NAME"

# install
mkdir -p "$INSTALL_DIR"
mv "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
chmod +x "$INSTALL_DIR/$BINARY_NAME"
ok "installed to $INSTALL_DIR/$BINARY_NAME"

# config dir
mkdir -p "$HOME/.config/ccc/profiles"

# check PATH
if ! echo "$PATH" | tr ':' '\n' | grep -qx "$INSTALL_DIR"; then
  warn "$INSTALL_DIR is not in your PATH"
  echo ""
  echo "  Add to your shell config:"
  if [[ -n "${ZSH_VERSION:-}" ]] || [[ "$SHELL" == */zsh ]]; then
    echo "    echo 'export PATH=\"$INSTALL_DIR:\$PATH\"' >> ~/.zshrc"
  else
    echo "    echo 'export PATH=\"$INSTALL_DIR:\$PATH\"' >> ~/.bashrc"
  fi
  echo ""
fi

ok "done! Run 'ccc' to start, or 'ccc --help' for usage."
