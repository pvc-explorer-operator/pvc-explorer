#!/usr/bin/env bash
set -euo pipefail

CLUSTER=pvc-explorer

read -r -p "Delete kind cluster '${CLUSTER}' and all data? [y/N] " confirm
[[ "$confirm" =~ ^[Yy]$ ]] || { echo "Aborted."; exit 0; }

kind delete cluster --name "$CLUSTER"
echo "Cluster deleted."
