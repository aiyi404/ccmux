#!/usr/bin/env bash
# ccc installer — downloads pre-built binary or builds from source
set -euo pipefail

GREEN='\033[32m' YELLOW='\033[33m' RED='\033[31m' RESET='\033[0m'
ok()   { printf "${GREEN}✓${RESET} %s\n" "$*"; }
warn() { printf "${YELLOW}warning:${RESET} %s\n" "$*"; }
die()  { printf "${RED}error:${RESET} %s\n" "$*" >&2; exit 1; }

INSTALL_DIR="${CCC_INSTALL_DIR:-$HOME/.local/bin}"
BINARY_NAME="ccc"
REPO="aiyi404/ccmux"

# Detect OS and architecture
detect_platform() {
  local os arch
  os="$(uname -s | tr '[:upper:]' '[:lower:]')"
  arch="$(uname -m)"

  case "$os" in
    linux)  os="linux" ;;
    darwin) os="darwin" ;;
    *)      die "Unsupported OS: $os" ;;
  esac

  case "$arch" in
    x86_64|amd64)   arch="amd64" ;;
    aarch64|arm64)   arch="arm64" ;;
    *)               die "Unsupported architecture: $arch" ;;
  esac

  echo "${os}_${arch}"
}

# Try downloading pre-built binary from GitHub Releases
try_download() {
  local platform="$1"
  local url="https://github.com/${REPO}/releases/latest/download/ccc_${platform}.tar.gz"
  local tmpdir
  tmpdir="$(mktemp -d)"

  echo "Downloading ccc for ${platform}..."
  if curl -fsSL "$url" -o "$tmpdir/ccc.tar.gz" 2>/dev/null; then
    tar -xzf "$tmpdir/ccc.tar.gz" -C "$tmpdir"
    if [[ -f "$tmpdir/ccc" ]]; then
      mkdir -p "$INSTALL_DIR"
      mv "$tmpdir/ccc" "$INSTALL_DIR/$BINARY_NAME"
      chmod +x "$INSTALL_DIR/$BINARY_NAME"
      rm -rf "$tmpdir"
      return 0
    fi
  fi
  rm -rf "$tmpdir"
  return 1
}

# Build from source
build_from_source() {
  command -v go &>/dev/null || die "No pre-built binary available and Go is not installed. Install Go from https://go.dev/dl/"

  local script_dir
  script_dir="$(cd "$(dirname "$0")" && pwd)"

  echo "Building from source..."
  cd "$script_dir"
  go build -ldflags="-s -w" -o "$BINARY_NAME" .
  mkdir -p "$INSTALL_DIR"
  mv "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
  chmod +x "$INSTALL_DIR/$BINARY_NAME"
}

# Main
platform="$(detect_platform)"

if try_download "$platform"; then
  ok "installed pre-built binary to $INSTALL_DIR/$BINARY_NAME"
else
  warn "no pre-built binary found, building from source..."
  build_from_source
  ok "built and installed to $INSTALL_DIR/$BINARY_NAME"
fi

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
