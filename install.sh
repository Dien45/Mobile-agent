#!/usr/bin/env sh
set -eu

GO_VERSION="1.24.2"
ROOT=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in x86_64|amd64) GOARCH=amd64;; aarch64|arm64) GOARCH=arm64;; *) echo "Unsupported architecture: $ARCH" >&2; exit 1;; esac

is_termux=false
if [ -n "${TERMUX_VERSION:-}" ] || [ -n "${PREFIX:-}" ] && printf '%s' "$PREFIX" | grep -q 'com.termux'; then is_termux=true; fi

if command -v go >/dev/null 2>&1; then
  GO=go
elif [ "$is_termux" = true ]; then
  echo "Installing Go from the Termux repository…"
  pkg install -y golang ca-certificates
  GO=go
else
  case "$OS" in linux|darwin) ;; *) echo "Unsupported Unix platform: $OS" >&2; exit 1;; esac
  CACHE="${XDG_CACHE_HOME:-$HOME/.cache}/arka/toolchains"
  mkdir -p "$CACHE"
  ARCHIVE="go${GO_VERSION}.${OS}-${GOARCH}.tar.gz"
  URL="https://go.dev/dl/$ARCHIVE"
  echo "Downloading pinned Go $GO_VERSION for $OS/$GOARCH…"
  curl -fL --retry 3 "$URL" -o "$CACHE/$ARCHIVE"
  EXPECTED=$(curl -fsSL "https://go.dev/dl/?mode=json&include=all" | awk -v archive="$ARCHIVE" '
    index($0, "\"filename\": \"" archive "\"") { found = 1 }
    found && /\"sha256\":/ {
      gsub(/[\",]/, "", $2)
      print $2
      exit
    }
  ')
  [ -n "$EXPECTED" ] || { echo "Could not read the Go checksum from official release metadata" >&2; exit 1; }
  ACTUAL=$(sha256sum "$CACHE/$ARCHIVE" 2>/dev/null | awk '{print $1}' || shasum -a 256 "$CACHE/$ARCHIVE" | awk '{print $1}')
  [ "$EXPECTED" = "$ACTUAL" ] || { echo "Go checksum verification failed" >&2; rm -f "$CACHE/$ARCHIVE"; exit 1; }
  rm -rf "$CACHE/go"
  tar -C "$CACHE" -xzf "$CACHE/$ARCHIVE"
  GO="$CACHE/go/bin/go"
fi

if [ -n "${ARKA_BIN_DIR:-}" ]; then
  PREFIX_DIR="$ARKA_BIN_DIR"
elif [ "$is_termux" = true ]; then
  PREFIX_DIR="${PREFIX}/bin"
elif [ -d /usr/local/bin ] && [ -w /usr/local/bin ]; then
  PREFIX_DIR="/usr/local/bin"
else
  PREFIX_DIR="$HOME/.local/bin"
fi
mkdir -p "$PREFIX_DIR"
echo "Building Arka locally…"
(cd "$ROOT" && "$GO" build -trimpath -ldflags="-s -w" -o "$PREFIX_DIR/arka" ./cmd/arka)
chmod 755 "$PREFIX_DIR/arka"

echo "Installed: $PREFIX_DIR/arka"
case ":$PATH:" in
  *":$PREFIX_DIR:"*) ;;
  *)
    case "${SHELL:-}" in
      */bash) PROFILE="$HOME/.bashrc" ;;
      */zsh) PROFILE="$HOME/.zshrc" ;;
      *) PROFILE="$HOME/.profile" ;;
    esac
    PATH_LINE="export PATH=\"$PREFIX_DIR:\$PATH\""
    touch "$PROFILE"
    grep -Fqx "$PATH_LINE" "$PROFILE" || printf '\n%s\n' "$PATH_LINE" >> "$PROFILE"
    echo "Added $PREFIX_DIR to PATH in $PROFILE"
    echo "Reload it now with: . $PROFILE"
    ;;
esac
echo "Run: arka start"
