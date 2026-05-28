#!/usr/bin/env bash
set -euo pipefail

CLUSTER=pvc-explorer
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# shellcheck source=kind/lib.sh
source "$SCRIPT_DIR/lib.sh"

ensure_tools kind

read -r -p "Delete kind cluster '${CLUSTER}' and all data? [y/N] " confirm
[[ "$confirm" =~ ^[Yy]$ ]] || { echo "Aborted."; exit 0; }

kind delete cluster --name "$CLUSTER"
echo "Cluster deleted."
