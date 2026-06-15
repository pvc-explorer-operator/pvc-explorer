#!/usr/bin/env bash
set -euo pipefail

CLUSTER=pvc-explorer
CONTROLLER_IMG=pvc-explorer:dev
AGENT_IMG=${AGENT_IMG:-ghcr.io/pvc-explorer-operator/pvc-explorer-agent:dev}
CONTROLLER_DIR="$(cd "$(dirname "$0")/.." && pwd)"
KIND_DIR="$CONTROLLER_DIR/kind"
SCRIPT_DIR="$KIND_DIR"

# shellcheck source=kind/lib.sh
source "$SCRIPT_DIR/lib.sh"

ensure_tools kind kubectl make go
detect_container_runtime

if kind get clusters 2>/dev/null | grep -q "^${CLUSTER}$"; then
  log "Cluster '${CLUSTER}' already exists — skipping creation"
else
  log "Creating kind cluster"
  kind create cluster --config "$KIND_DIR/cluster.yaml"
fi

kubectl cluster-info --context "kind-${CLUSTER}" >/dev/null

log "Building controller image"
$DOCKER build -t "$CONTROLLER_IMG" "$CONTROLLER_DIR"

log "Pulling agent image"
if ! $DOCKER pull "$AGENT_IMG"; then
  echo "Could not pull agent image: $AGENT_IMG" >&2
  echo "Note: GHCR package visibility is separate from repository visibility" >&2
  echo "Try one of:" >&2
  echo "  1) $DOCKER login ghcr.io" >&2
  echo "  2) Use a different image: AGENT_IMG=<image> kind/setup.sh" >&2
  exit 1
fi

log "Loading images into kind"
kind_load "$CONTROLLER_IMG"
kind_load "$AGENT_IMG"

log "Installing CRDs"
make -C "$CONTROLLER_DIR" manifests
kubectl apply -k "$CONTROLLER_DIR/config/crd"


log "Creating system namespace"
kubectl create namespace pvc-explorer-system --dry-run=client -o yaml | kubectl apply -f -

log "Deploying controller (dev overlay)"
kubectl apply -k "$CONTROLLER_DIR/kustomize/overlays/dev"

log "Restarting controller to pick up rebuilt image"
kubectl rollout restart deployment -n pvc-explorer-system pvc-explorer-controller-manager

log "Waiting for controller rollout"
kubectl rollout status deployment -n pvc-explorer-system \
  pvc-explorer-controller-manager --timeout=300s

# Password hash is bcrypt of "admin". Regenerate with:
#   htpasswd -nbB admin admin | cut -d: -f2
BCRYPT_ADMIN='$2a$10$UAvOdwl6OeMNAlXbSVwKH.ag86u60RkCUDMTsQHdTnpO7o/msX6SK'

log "Creating auth secret (admin/admin)"
kubectl delete secret pvc-explorer-auth -n pvc-explorer-system --ignore-not-found
kubectl create secret generic pvc-explorer-auth \
  -n pvc-explorer-system \
  --from-literal="admin=${BCRYPT_ADMIN}"

log "Creating config maps (pvc-explorer-config, pvc-explorer-rbac)"
kubectl apply -f - <<'EOF'
apiVersion: v1
kind: ConfigMap
metadata:
  name: pvc-explorer-config
  namespace: pvc-explorer-system
data:
  adminUsers: "admin"
  oidc.enabled: "false"
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: pvc-explorer-rbac
  namespace: pvc-explorer-system
data:
  policy.default: "viewer"
  policy.csv: ""
EOF

# ── Dex OIDC provider (local, no external connections) ─────────────────────
log "Deploying Dex OIDC provider"

# Generate self-signed TLS cert for Dex
DEX_CERT_DIR=$(mktemp -d)
trap 'rm -rf "$DEX_CERT_DIR"' EXIT

openssl req -x509 -nodes -newkey rsa:2048 \
  -keyout "$DEX_CERT_DIR/tls.key" \
  -out "$DEX_CERT_DIR/tls.crt" \
  -days 3650 \
  -subj "/CN=dex.pvc-explorer-system.svc.cluster.local" \
  -addext "subjectAltName=DNS:dex.pvc-explorer-system.svc.cluster.local,DNS:localhost,IP:127.0.0.1" \
  2>/dev/null

# Create Dex TLS secret
kubectl delete secret dex-tls -n pvc-explorer-system --ignore-not-found
kubectl create secret tls dex-tls \
  -n pvc-explorer-system \
  --cert="$DEX_CERT_DIR/tls.crt" \
  --key="$DEX_CERT_DIR/tls.key"

# Pre-computed bcrypt hashes (cost=10):
#   admin123  -> $2a$10$LvpSVk0ccCuniH3hMUBKgOTXWOkRWL5VXqwH/5RY8NwN1fhm8w.6.
#   user123   -> $2a$10$7i1aaEfcMwO0jZeqBhd3uuoh8xTZI5bT/3XctydgVYPm3Rkpd42.i
#   viewer123 -> $2a$10$YmQt5lkYoxoIAeT0tHC8qu/e4e94ivQqMKSOMMbn92EJe4w5uS99S
DEX_ADMIN_HASH='$2a$10$LvpSVk0ccCuniH3hMUBKgOTXWOkRWL5VXqwH/5RY8NwN1fhm8w.6.'
DEX_USER_HASH='$2a$10$7i1aaEfcMwO0jZeqBhd3uuoh8xTZI5bT/3XctydgVYPm3Rkpd42.i'
DEX_VIEWER_HASH='$2a$10$YmQt5lkYoxoIAeT0tHC8qu/e4e94ivQqMKSOMMbn92EJe4w5uS99S'

kubectl apply -f - <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: dex-config
  namespace: pvc-explorer-system
data:
  config.yaml: |
    issuer: https://dex.pvc-explorer-system.svc.cluster.local:5556
    storage:
      type: memory
    web:
      https: 0.0.0.0:5556
      tlsCert: /etc/dex/tls/tls.crt
      tlsKey: /etc/dex/tls/tls.key
    staticPasswords:
      - email: admin@pvc-explorer.local
        username: admin
        hash: ${DEX_ADMIN_HASH}
        groups:
          - platform-admins
      - email: user@pvc-explorer.local
        username: user
        hash: ${DEX_USER_HASH}
        groups:
          - platform-users
      - email: viewer@pvc-explorer.local
        username: viewer
        hash: ${DEX_VIEWER_HASH}
        groups:
          - platform-viewers
    staticClients:
      - id: pvc-explorer
        redirectURIs:
          - 'http://localhost:8080/api/v1/auth/oidc/callback'
        name: 'PVC Explorer'
        secret: pvc-explorer-dex-secret
    enablePasswordDB: true
    oauth2:
      responseTypes: ['code']
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dex
  namespace: pvc-explorer-system
  labels:
    app: dex
spec:
  replicas: 1
  selector:
    matchLabels:
      app: dex
  template:
    metadata:
      labels:
        app: dex
    spec:
      containers:
        - name: dex
          image: ghcr.io/dexidp/dex:v2.45.1
          command: ["dex", "serve", "/etc/dex/config.yaml"]
          ports:
            - containerPort: 5556
              name: https
          volumeMounts:
            - name: config
              mountPath: /etc/dex/config.yaml
              subPath: config.yaml
            - name: tls
              mountPath: /etc/dex/tls
      volumes:
        - name: config
          configMap:
            name: dex-config
        - name: tls
          secret:
            secretName: dex-tls
---
apiVersion: v1
kind: Service
metadata:
  name: dex
  namespace: pvc-explorer-system
spec:
  type: NodePort
  selector:
    app: dex
  ports:
    - port: 5556
      targetPort: 5556
      nodePort: 30556
      protocol: TCP
      name: https
EOF

log "Waiting for Dex to be ready"
kubectl rollout status deployment/dex -n pvc-explorer-system --timeout=60s

log "Configuring pvc-explorer OIDC (Dex)"
kubectl apply -f - <<'EOF'
apiVersion: v1
kind: ConfigMap
metadata:
  name: pvc-explorer-config
  namespace: pvc-explorer-system
data:
  adminUsers: "admin"
  oidc.enabled: "true"
  oidc.issuer: "https://dex.pvc-explorer-system.svc.cluster.local:5556"
  oidc.externalIssuer: "https://localhost:5556"
  oidc.clientID: "pvc-explorer"
  oidc.clientSecret: "$pvc-explorer-oidc.clientSecret"
  oidc.redirectURI: "http://localhost:8080/api/v1/auth/oidc/callback"
  oidc.scopes: "openid,profile,email,groups"
  oidc.groupClaim: "groups"
  oidc.skipTLSVerify: "true"
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: pvc-explorer-rbac
  namespace: pvc-explorer-system
data:
  policy.default: "viewer"
  policy.csv: |
    g, "platform-admins", admin
    g, "platform-users", user
    g, "platform-viewers", viewer
EOF

log "Creating OIDC client secret"
kubectl delete secret pvc-explorer-oidc -n pvc-explorer-system --ignore-not-found
kubectl create secret generic pvc-explorer-oidc \
  -n pvc-explorer-system \
  --from-literal=clientSecret='pvc-explorer-dex-secret'

log "Restarting controller to pick up OIDC config"
kubectl rollout restart deployment -n pvc-explorer-system pvc-explorer-controller-manager
kubectl rollout status deployment -n pvc-explorer-system \
  pvc-explorer-controller-manager --timeout=300s

log "Applying StorageClass (Immediate binding)"
kubectl apply -f "$KIND_DIR/storageclass-immediate.yaml"
kubectl patch storageclass standard \
  -p '{"metadata":{"annotations":{"storageclass.kubernetes.io/is-default-class":"false"}}}' \
  2>/dev/null || true

NODE=$(kubectl get nodes -o jsonpath='{.items[0].metadata.name}')

log "Pre-creating hostPath PVs for static demo PVCs"
for entry in demo-data:1Gi demo-logs:512Mi demo-cache:256Mi demo-unclaimed:1Gi; do
  name=${entry%%:*}
  size=${entry##*:}
  kubectl apply -f - <<EOF
apiVersion: v1
kind: PersistentVolume
metadata:
  name: ${name}
spec:
  capacity:
    storage: ${size}
  accessModes:
    - ReadWriteOnce
  persistentVolumeReclaimPolicy: Delete
  storageClassName: demo-hostpath
  hostPath:
    path: /tmp/pvc-explorer/${name}
    type: DirectoryOrCreate
  nodeAffinity:
    required:
      nodeSelectorTerms:
        - matchExpressions:
            - key: kubernetes.io/hostname
              operator: In
              values:
                - ${NODE}
EOF
done

log "Pre-creating hostPath PVs for demo-staging PVCs"
for entry in staging-app:2Gi staging-db:5Gi staging-uploads:1Gi; do
  name=${entry%%:*}
  size=${entry##*:}
  mode=ReadWriteOnce
  [[ "$name" == "staging-uploads" ]] && mode=ReadWriteMany
  kubectl apply -f - <<EOF
apiVersion: v1
kind: PersistentVolume
metadata:
  name: ${name}
spec:
  capacity:
    storage: ${size}
  accessModes:
    - ${mode}
  persistentVolumeReclaimPolicy: Delete
  storageClassName: demo-hostpath
  hostPath:
    path: /tmp/pvc-explorer/${name}
    type: DirectoryOrCreate
  nodeAffinity:
    required:
      nodeSelectorTerms:
        - matchExpressions:
            - key: kubernetes.io/hostname
              operator: In
              values:
                - ${NODE}
EOF
done

log "Creating demo namespace and PVCs"
kubectl apply -f "$KIND_DIR/demo-namespace.yaml"
kubectl apply -f "$KIND_DIR/demo-pvcs.yaml"
kubectl apply -f "$KIND_DIR/demo-pvc-generator.yaml"

log "Creating demo-staging namespace and PVCs"
kubectl apply -f "$KIND_DIR/demo-staging-namespace.yaml"
kubectl apply -f "$KIND_DIR/demo-staging-pvcs.yaml"

log "Applying demo PVCExplorerScope (explicit namespace names)"
kubectl apply -f "$KIND_DIR/demo-scope.yaml"

log "Applying demo-by-label PVCExplorerScope (labelSelector)"
kubectl apply -f "$KIND_DIR/demo-label-scope.yaml"

echo ""
echo "────────────────────────────────────────────"
echo " Cluster ready: kind-${CLUSTER}"
echo ""
echo " Dashboard:  http://localhost:8080  (admin / admin)"
echo ""
echo " OIDC (Dex): https://localhost:5556"
echo "   Test users (password in parentheses):"
echo "     admin@pvc-explorer.local   (admin123)   → role: admin"
echo "     user@pvc-explorer.local    (user123)    → role: user"
echo "     viewer@pvc-explorer.local  (viewer123)  → role: viewer"
echo ""
echo " Namespaces:"
echo "   demo           — managed by scope 'demo' (explicit names)"
echo "   demo-staging   — managed by scope 'demo-by-label' (labelSelector)"
echo ""
echo " Watch CRs:"
echo "   kubectl get pvcexplorerscope,pvcexplorer -A"
echo ""
echo " Controller logs:"
echo "   kubectl logs -n pvc-explorer-system"
echo "     -l control-plane=controller-manager -f"
echo "────────────────────────────────────────────"
