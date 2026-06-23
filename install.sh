#!/bin/sh
# apikit installer (macOS / Linux).
#
# Downloads a release archive from GitHub, verifies its SHA-256 against
# checksums.txt, and installs the `apikit` binary into a directory on PATH.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/dkryvak/apikit/main/install.sh | sh
#   ./install.sh [-v VERSION] [-d DIR] [--no-modify-path]
#
# Options / environment:
#   -v VERSION       | APIKIT_VERSION         release to install (default: latest)
#   -d DIR           | APIKIT_INSTALL_DIR     install dir (default: /usr/local/bin,
#                                             via sudo when not writable)
#   --no-modify-path | APIKIT_NO_MODIFY_PATH  do not edit shell rc files for PATH
set -eu

REPO_OWNER="dkryvak"
REPO_NAME="apikit"
BIN="apikit"

VERSION="${APIKIT_VERSION:-}"
INSTALL_DIR="${APIKIT_INSTALL_DIR:-}"
NO_MODIFY_PATH="${APIKIT_NO_MODIFY_PATH:-}"

while [ $# -gt 0 ]; do
  case "$1" in
    -v) VERSION="$2"; shift 2 ;;
    -d) INSTALL_DIR="$2"; shift 2 ;;
    --no-modify-path) NO_MODIFY_PATH=1; shift ;;
    -h|--help) sed -n '2,21p' "$0" | sed 's/^# \{0,1\}//'; exit 0 ;;
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
# Default to /usr/local/bin (already on PATH on macOS/Linux). When it isn't
# writable, the install step below uses sudo instead of silently falling back to
# a per-user dir. A custom dir is honored only via -d / APIKIT_INSTALL_DIR.
if [ -z "$INSTALL_DIR" ]; then INSTALL_DIR="/usr/local/bin"; fi

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

# --- PATH setup ------------------------------------------------------------
# If the install dir isn't on PATH (typically only with a custom -d), append an
# idempotent export to the right rc file for the user's login shell. The block
# is guarded at runtime so sourcing it twice can't duplicate the entry, and we
# skip writing when an identical line already exists. Opt out: --no-modify-path.
add_to_path() {
  dir="$1"

  # Already reachable — nothing to do.
  case ":${PATH}:" in *":${dir}:"*) return 0 ;; esac

  guard="case \":\$PATH:\" in *\":${dir}:\"*) ;; *) export PATH=\"${dir}:\$PATH\" ;; esac"

  if [ -n "$NO_MODIFY_PATH" ]; then
    info "${dir} is not on your PATH. Add it manually:"
    echo "      export PATH=\"${dir}:\$PATH\"" >&2
    return 0
  fi

  # rc file(s) to update, chosen by the login shell.
  case "$(basename "${SHELL:-/bin/sh}")" in
    zsh)  rcs="${ZDOTDIR:-$HOME}/.zshrc" ;;
    bash) rcs="${HOME}/.bash_profile ${HOME}/.bashrc" ;;
    *)    rcs="${HOME}/.profile" ;;
  esac

  updated=""
  for rc in $rcs; do
    [ -e "$rc" ] || : > "$rc"
    if ! grep -qF "$guard" "$rc" 2>/dev/null; then
      printf '\n# added by apikit installer\n%s\n' "$guard" >> "$rc"
      updated="${updated} ${rc}"
    fi
  done

  if [ -n "$updated" ]; then
    info "added ${dir} to your PATH in:${updated}"
    set -- $updated
    info "restart your shell or run: source $1"
  fi
}

add_to_path "$INSTALL_DIR"

"$dest" --version 2>/dev/null || true
