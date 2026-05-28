#!/usr/bin/env bash
set -euo pipefail

CLUSTER=pvc-explorer
CONTROLLER_IMG=pvc-explorer:dev
AGENT_IMG=${AGENT_IMG:-ghcr.io/pvc-explorer-operator/pvc-explorer-agent:dev}
CONTROLLER_DIR="$(cd "$(dirname "$0")/.." && pwd)"

log() { echo "==> $*"; }

# kind load docker-image only works with Docker's image store.
# For Podman, export the image to a tarball and use image-archive.
kind_load() {
  local img="$1"
  if [ "$DOCKER" = "podman" ]; then
    local tmp
    tmp=$(mktemp /tmp/kind-img-XXXXXX.tar)
    podman save "$img" -o "$tmp"
    kind load image-archive "$tmp" --name "$CLUSTER"
    rm -f "$tmp"
  else
    kind load docker-image "$img" --name "$CLUSTER"
  fi
}

if command -v docker >/dev/null 2>&1; then
  DOCKER=docker
elif command -v podman >/dev/null 2>&1; then
  DOCKER=podman
  log "docker not found — using podman"
else
  echo "Neither docker nor podman found" >&2; exit 1
fi

TARGET="${1:-both}"

if [[ "$TARGET" == "controller" || "$TARGET" == "both" ]]; then
  log "Rebuilding controller"
  $DOCKER build -t "$CONTROLLER_IMG" "$CONTROLLER_DIR"
  SHORT_SHA=$($DOCKER inspect --format='{{.Id}}' "$CONTROLLER_IMG" | cut -c8-19)
  SHA_TAG="pvc-explorer:${SHORT_SHA}"
  $DOCKER tag "$CONTROLLER_IMG" "$SHA_TAG"
  kind_load "$SHA_TAG"
  kubectl set image deployment \
    -n pvc-explorer-system \
    pvc-explorer-controller-manager \
    "manager=$SHA_TAG"
  kubectl rollout status deployment -n pvc-explorer-system \
    pvc-explorer-controller-manager --timeout=60s
fi

if [[ "$TARGET" == "agent" || "$TARGET" == "both" ]]; then
  log "Pulling agent image"
  if ! $DOCKER pull "$AGENT_IMG"; then
    echo "Could not pull agent image: $AGENT_IMG" >&2
    echo "Note: GHCR package visibility is separate from repository visibility" >&2
    echo "Try one of:" >&2
    echo "  1) $DOCKER login ghcr.io" >&2
    echo "  2) Use a different image: AGENT_IMG=<image> kind/rebuild.sh agent" >&2
    exit 1
  fi
  kind_load "$AGENT_IMG"
  log "Agent image reloaded — existing agent pods will pick it up on next wake cycle"
fi

log "Done"
