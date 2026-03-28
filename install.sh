#!/usr/bin/env bash
# ccc installer — copies the script to ~/.local/bin and sets up config dir
set -euo pipefail

BOLD='\033[1m' GREEN='\033[32m' YELLOW='\033[33m' RED='\033[31m' RESET='\033[0m'
ok()   { printf "${GREEN}✓${RESET} %s\n" "$*"; }
warn() { printf "${YELLOW}warning:${RESET} %s\n" "$*"; }
die()  { printf "${RED}error:${RESET} %s\n" "$*" >&2; exit 1; }

INSTALL_DIR="${CCC_INSTALL_DIR:-$HOME/.local/bin}"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
SOURCE="$SCRIPT_DIR/ccc"

[[ -f "$SOURCE" ]] || die "ccc script not found at $SOURCE"

# check jq
if ! command -v jq &>/dev/null; then
  warn "jq is not installed. ccc requires jq to run."
  echo "  macOS:  brew install jq"
  echo "  Ubuntu: sudo apt install jq"
  echo "  Arch:   sudo pacman -S jq"
fi

# install
mkdir -p "$INSTALL_DIR"
cp "$SOURCE" "$INSTALL_DIR/ccc"
chmod +x "$INSTALL_DIR/ccc"
ok "installed ccc to $INSTALL_DIR/ccc"

# config dir
mkdir -p "$HOME/.config/ccc/profiles"
ok "created config directory ~/.config/ccc/"

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

ok "done! Run 'ccc --help' to get started."
