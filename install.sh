#!/bin/sh
# apikit installer (macOS / Linux).
#
# Downloads a release archive from GitHub, verifies its SHA-256 against
# checksums.txt, and installs the `apikit` binary into a directory on PATH.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/dkryvak/apikit/main/install.sh | sh
#   ./install.sh [-v VERSION] [-d DIR]
#
# Options / environment:
#   -v VERSION  | APIKIT_VERSION       release to install, e.g. v0.2.0 (default: latest)
#   -d DIR      | APIKIT_INSTALL_DIR   install directory (default: /usr/local/bin if
#                                      writable, else ~/.local/bin)
set -eu

REPO_OWNER="dkryvak"
REPO_NAME="apikit"
BIN="apikit"

VERSION="${APIKIT_VERSION:-}"
INSTALL_DIR="${APIKIT_INSTALL_DIR:-}"

while [ $# -gt 0 ]; do
  case "$1" in
    -v) VERSION="$2"; shift 2 ;;
    -d) INSTALL_DIR="$2"; shift 2 ;;
    -h|--help) sed -n '2,20p' "$0" | sed 's/^# \{0,1\}//'; exit 0 ;;
    *) echo "unknown option: $1" >&2; exit 2 ;;
  esac
done

err() { echo "error: $*" >&2; exit 1; }
info() { echo "==> $*" >&2; }

have() { command -v "$1" >/dev/null 2>&1; }

# --- prerequisites ---------------------------------------------------------
if have curl; then DL="curl -fsSL"; DLO="curl -fsSL -o";
elif have wget; then DL="wget -qO-"; DLO="wget -qO";
else err "need curl or wget"; fi

have tar || err "need tar"

if have sha256sum; then SHACMD="sha256sum";
elif have shasum; then SHACMD="shasum -a 256";
else err "need sha256sum or shasum"; fi

# --- detect OS / arch ------------------------------------------------------
os="$(uname -s)"
case "$os" in
  Darwin) OS="darwin" ;;
  Linux)  OS="linux" ;;
  *) err "unsupported OS: $os" ;;
esac

arch="$(uname -m)"
case "$arch" in
  x86_64|amd64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) err "unsupported architecture: $arch" ;;
esac

# --- resolve version -------------------------------------------------------
if [ -z "$VERSION" ]; then
  info "resolving latest release"
  api="https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest"
  VERSION="$($DL "$api" | grep -m1 '"tag_name"' | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')"
  [ -n "$VERSION" ] || err "could not determine latest version"
fi
case "$VERSION" in v*) ;; *) VERSION="v${VERSION}" ;; esac
# GoReleaser strips the leading 'v' for asset names (.Version).
VER_NUM="${VERSION#v}"

ASSET="${BIN}_${VER_NUM}_${OS}_${ARCH}.tar.gz"
BASE="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${VERSION}"

# --- choose install dir ----------------------------------------------------
if [ -z "$INSTALL_DIR" ]; then
  if [ -w /usr/local/bin ] 2>/dev/null; then INSTALL_DIR="/usr/local/bin";
  else INSTALL_DIR="${HOME}/.local/bin"; fi
fi

# --- download + verify -----------------------------------------------------
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

info "downloading ${ASSET} (${VERSION})"
$DLO "${TMP}/${ASSET}" "${BASE}/${ASSET}" || err "download failed: ${BASE}/${ASSET}"
$DLO "${TMP}/checksums.txt" "${BASE}/checksums.txt" || err "download failed: checksums.txt"

info "verifying sha256"
expected="$(grep " ${ASSET}\$" "${TMP}/checksums.txt" | awk '{print $1}')"
[ -n "$expected" ] || err "no checksum entry for ${ASSET}"
actual="$(cd "$TMP" && $SHACMD "$ASSET" | awk '{print $1}')"
[ "$expected" = "$actual" ] || err "checksum mismatch (expected $expected, got $actual)"

# --- extract + install -----------------------------------------------------
tar -xzf "${TMP}/${ASSET}" -C "$TMP" "$BIN" || err "could not extract ${BIN}"

mkdir -p "$INSTALL_DIR"
dest="${INSTALL_DIR}/${BIN}"
if [ -w "$INSTALL_DIR" ]; then
  install -m 0755 "${TMP}/${BIN}" "$dest"
else
  info "$INSTALL_DIR is not writable, using sudo"
  sudo install -m 0755 "${TMP}/${BIN}" "$dest" || err "install to $dest failed"
fi

# macOS: clear the quarantine flag so Gatekeeper doesn't block the binary.
if [ "$OS" = "darwin" ] && have xattr; then
  xattr -d com.apple.quarantine "$dest" 2>/dev/null || true
fi

info "installed ${BIN} ${VERSION} -> ${dest}"

# --- PATH hint -------------------------------------------------------------
case ":${PATH}:" in
  *":${INSTALL_DIR}:"*) ;;
  *) info "note: ${INSTALL_DIR} is not in your PATH — add it, e.g.:"
     echo "      export PATH=\"${INSTALL_DIR}:\$PATH\"" >&2 ;;
esac

"$dest" --version 2>/dev/null || true
