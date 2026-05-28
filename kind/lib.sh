#!/usr/bin/env bash
# lib.sh — shared helpers for kind/ scripts
# Sourced by setup.sh, rebuild.sh, teardown.sh

log()  { echo "==> $*"; }
info() { echo "    $*"; }

# ── OS / package-manager detection ──────────────────────────────────────────
detect_os() {
  if [[ "${OSTYPE:-}" == darwin* ]]; then
    OS=macos
  elif grep -qi microsoft /proc/version 2>/dev/null; then
    OS=wsl
    # WSL: detect deb vs yum inside
    if [ -f /etc/debian_version ]; then
      PKG=deb
    elif [ -f /etc/redhat-release ] || [ -f /etc/fedora-release ]; then
      PKG=yum
    else
      PKG=unknown
    fi
  elif [ -f /etc/debian_version ]; then
    OS=linux; PKG=deb
  elif [ -f /etc/redhat-release ] || [ -f /etc/fedora-release ]; then
    OS=linux; PKG=yum
  else
    OS=unknown; PKG=unknown
  fi
  [ "${OS:-}" = "macos" ] && PKG=brew
  [ "${OS:-}" = "wsl" ]   && OS=wsl
}

_arch() {
  local m; m=$(uname -m)
  case "$m" in
    x86_64)  echo amd64 ;;
    aarch64|arm64) echo arm64 ;;
    *) echo "$m" ;;
  esac
}

# ── Per-tool installers ───────────────────────────────────────────────────────
_install_brew_or_fail() {
  local tool="$1"
  if command -v brew >/dev/null 2>&1; then
    brew install "$tool"
  else
    echo "Homebrew not found. Install it from https://brew.sh then re-run." >&2
    exit 1
  fi
}

_install_kind() {
  log "Installing kind…"
  case "$PKG" in
    brew) _install_brew_or_fail kind ;;
    deb|yum)
      local arch; arch=$(_arch)
      curl -sSLo /tmp/kind \
        "https://kind.sigs.k8s.io/dl/latest/kind-linux-${arch}"
      chmod +x /tmp/kind
      sudo mv /tmp/kind /usr/local/bin/kind
      ;;
    *) echo "Cannot auto-install kind on this OS. See https://kind.sigs.k8s.io" >&2; exit 1 ;;
  esac
}

_install_kubectl() {
  log "Installing kubectl…"
  case "$PKG" in
    brew) _install_brew_or_fail kubectl ;;
    deb|yum)
      local arch; arch=$(_arch)
      local ver; ver=$(curl -sSL https://dl.k8s.io/release/stable.txt)
      curl -sSLo /tmp/kubectl \
        "https://dl.k8s.io/release/${ver}/bin/linux/${arch}/kubectl"
      chmod +x /tmp/kubectl
      sudo mv /tmp/kubectl /usr/local/bin/kubectl
      ;;
    *) echo "Cannot auto-install kubectl on this OS. See https://kubernetes.io/docs/tasks/tools" >&2; exit 1 ;;
  esac
}

_install_make() {
  log "Installing make…"
  case "$PKG" in
    brew) xcode-select --install 2>/dev/null || _install_brew_or_fail make ;;
    deb)  sudo apt-get install -y --no-install-recommends make ;;
    yum)
      if command -v dnf >/dev/null 2>&1; then sudo dnf install -y make
      else sudo yum install -y make; fi
      ;;
    *) echo "Cannot auto-install make on this OS." >&2; exit 1 ;;
  esac
}

_install_go() {
  log "Installing Go…"
  case "$PKG" in
    brew) _install_brew_or_fail go ;;
    deb)  sudo apt-get install -y --no-install-recommends golang-go ;;
    yum)
      if command -v dnf >/dev/null 2>&1; then sudo dnf install -y golang
      else sudo yum install -y golang; fi
      ;;
    *)
      echo "Cannot auto-install Go on this OS. See https://go.dev/doc/install" >&2
      exit 1
      ;;
  esac
}

_install_docker() {
  log "Installing Docker…"
  case "$PKG" in
    brew)
      echo "Install Docker Desktop for Mac from https://docs.docker.com/desktop/install/mac-install/" >&2
      exit 1
      ;;
    deb)
      sudo apt-get update -qq
      sudo apt-get install -y --no-install-recommends \
        ca-certificates curl gnupg lsb-release
      curl -fsSL https://download.docker.com/linux/ubuntu/gpg \
        | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
      echo \
        "deb [arch=$(_arch) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] \
        https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" \
        | sudo tee /etc/apt/sources.list.d/docker.list >/dev/null
      sudo apt-get update -qq
      sudo apt-get install -y docker-ce docker-ce-cli containerd.io
      ;;
    yum)
      sudo yum install -y yum-utils
      sudo yum-config-manager --add-repo \
        https://download.docker.com/linux/centos/docker-ce.repo
      sudo yum install -y docker-ce docker-ce-cli containerd.io
      sudo systemctl enable --now docker
      ;;
    *) echo "Cannot auto-install Docker on this OS. See https://docs.docker.com/engine/install" >&2; exit 1 ;;
  esac
}

# ── Container runtime detection ───────────────────────────────────────────────
detect_container_runtime() {
  if command -v docker >/dev/null 2>&1; then
    DOCKER=docker
  elif command -v podman >/dev/null 2>&1; then
    DOCKER=podman
    log "docker not found — using podman"
  else
    log "No container runtime found. Attempting to install Docker…"
    _install_docker
    DOCKER=docker
  fi
}

# ── kind image loader (Docker store vs Podman store) ─────────────────────────
kind_load() {
  local img="$1"
  if [ "$DOCKER" = "podman" ]; then
    local tmp
    tmp=$(mktemp /tmp/kind-img-XXXXXX.tar)

    # Podman stores images as localhost/<name>. containerd (inside kind)
    # resolves bare image names (e.g. pvc-explorer:dev) as
    # docker.io/library/pvc-explorer:dev, so it can't find the locally
    # loaded image and falls back to a registry pull (auth failure).
    #
    # Fix: strip any existing registry prefix, then re-tag with
    # docker.io/library/ before saving so containerd finds the image locally.
    local bare="${img##*/}"   # strip registry/org prefix if any
    local full="docker.io/library/${bare}"
    podman tag "$img" "$full"
    podman save --format docker-archive "$full" -o "$tmp"
    podman rmi "$full" >/dev/null 2>&1 || true

    kind load image-archive "$tmp" --name "$CLUSTER"
    rm -f "$tmp"
  else
    kind load docker-image "$img" --name "$CLUSTER"
  fi
}

# ── Tool checker / installer ──────────────────────────────────────────────────
# Usage: ensure_tools kind kubectl make go
ensure_tools() {
  detect_os
  local missing=()
  for tool in "$@"; do
    command -v "$tool" >/dev/null 2>&1 || missing+=("$tool")
  done

  if [ ${#missing[@]} -eq 0 ]; then
    return 0
  fi

  log "Missing tools: ${missing[*]}"
  for tool in "${missing[@]}"; do
    case "$tool" in
      kind)    _install_kind    ;;
      kubectl) _install_kubectl ;;
      make)    _install_make    ;;
      go)      _install_go      ;;
      docker)  _install_docker  ;;
      *)
        echo "No auto-installer for '$tool'. Please install it manually." >&2
        exit 1
        ;;
    esac
  done

  # Verify all tools are now available
  for tool in "$@"; do
    command -v "$tool" >/dev/null 2>&1 || {
      echo "Installation of '$tool' failed or not in PATH." >&2
      exit 1
    }
  done
  log "All required tools are available"
}
