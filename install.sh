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

# check gum (optional, for TUI)
if command -v gum &>/dev/null; then
  ok "gum found ($(gum --version))"
else
  echo ""
  echo "${BOLD}gum${RESET} is not installed. Install it for interactive TUI support."
  echo ""
  # detect OS and offer install
  install_gum=""
  if command -v brew &>/dev/null; then
    install_gum="brew install gum"
  elif command -v pacman &>/dev/null; then
    install_gum="sudo pacman -S gum"
  elif command -v apt-get &>/dev/null; then
    install_gum="apt-get"
  elif command -v dnf &>/dev/null; then
    install_gum="sudo dnf install gum"
  elif command -v go &>/dev/null; then
    install_gum="go install github.com/charmbracelet/gum@latest"
  fi

  if [[ -n "$install_gum" ]]; then
    read -rp "Install gum now? [y/N] " ans
    if [[ "$ans" =~ ^[yY]$ ]]; then
      if [[ "$install_gum" == "apt-get" ]]; then
        # charm apt repo
        echo "  Adding Charm apt repository..."
        sudo mkdir -p /etc/apt/keyrings
        curl -fsSL https://repo.charm.sh/apt/gpg.key | sudo gpg --dearmor -o /etc/apt/keyrings/charm.gpg 2>/dev/null
        echo "deb [signed-by=/etc/apt/keyrings/charm.gpg] https://repo.charm.sh/apt/ * *" | sudo tee /etc/apt/sources.list.d/charm.list > /dev/null
        sudo apt-get update -qq && sudo apt-get install -y -qq gum
      elif [[ "$install_gum" == go* ]]; then
        echo "  Installing via go..."
        $install_gum
        # link to install dir if go bin is not in PATH
        if ! command -v gum &>/dev/null && [[ -f "$HOME/go/bin/gum" ]]; then
          ln -sf "$HOME/go/bin/gum" "$INSTALL_DIR/gum"
        fi
      else
        echo "  Running: $install_gum"
        $install_gum
      fi
      if command -v gum &>/dev/null || [[ -f "$INSTALL_DIR/gum" ]]; then
        ok "gum installed"
      else
        warn "gum installation failed. You can install it manually later."
      fi
    else
      echo "  Skip. You can install gum later for TUI support:"
      echo "    macOS:  brew install gum"
      echo "    Ubuntu: see https://github.com/charmbracelet/gum#installation"
      echo "    Go:     go install github.com/charmbracelet/gum@latest"
    fi
  else
    echo "  Install manually: https://github.com/charmbracelet/gum#installation"
  fi
  echo ""
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
