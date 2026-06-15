#!/usr/bin/env bash
# kind/dex/generate-cert.sh — generates a self-signed TLS cert for Dex
# Called by setup.sh. Outputs cert and key to stdout (base64-encoded).
set -euo pipefail

TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

openssl req -x509 -nodes -newkey rsa:2048 \
  -keyout "$TMPDIR/tls.key" \
  -out "$TMPDIR/tls.crt" \
  -days 3650 \
  -subj "/CN=dex.pvc-explorer-system.svc.cluster.local" \
  -addext "subjectAltName=DNS:dex.pvc-explorer-system.svc.cluster.local,DNS:localhost,IP:127.0.0.1" \
  2>/dev/null

cat "$TMPDIR/tls.crt"
echo "---KEY---"
cat "$TMPDIR/tls.key"
