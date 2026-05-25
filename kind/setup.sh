#!/usr/bin/env bash
set -euo pipefail

CLUSTER=pvc-explorer
CONTROLLER_IMG=pvc-explorer:dev
AGENT_IMG=${AGENT_IMG:-ghcr.io/pvc-explorer-operator/pvc-explorer-agent:dev}
CONTROLLER_DIR="$(cd "$(dirname "$0")/.." && pwd)"
KIND_DIR="$CONTROLLER_DIR/kind"

log() { echo "==> $*"; }

need() {
  for cmd in "$@"; do
    command -v "$cmd" >/dev/null 2>&1 || { echo "Missing: $cmd" >&2; exit 1; }
  done
}

need kind kubectl docker

if kind get clusters 2>/dev/null | grep -q "^${CLUSTER}$"; then
  log "Cluster '${CLUSTER}' already exists — skipping creation"
else
  log "Creating kind cluster"
  kind create cluster --config "$KIND_DIR/cluster.yaml"
fi

kubectl cluster-info --context "kind-${CLUSTER}" >/dev/null

log "Building controller image"
docker build -t "$CONTROLLER_IMG" "$CONTROLLER_DIR"

log "Pulling agent image"
if ! docker pull "$AGENT_IMG"; then
  echo "Could not pull agent image: $AGENT_IMG" >&2
  echo "Note: GHCR package visibility is separate from repository visibility" >&2
  echo "Try one of:" >&2
  echo "  1) docker login ghcr.io" >&2
  echo "  2) Use a different image: AGENT_IMG=<image> kind/setup.sh" >&2
  exit 1
fi

log "Loading images into kind"
kind load docker-image "$CONTROLLER_IMG" --name "$CLUSTER"
kind load docker-image "$AGENT_IMG"      --name "$CLUSTER"

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
  pvc-explorer-controller-manager --timeout=120s

# Password hash is bcrypt of "admin". Regenerate with:
#   htpasswd -nbB admin admin | cut -d: -f2
BCRYPT_ADMIN='$2a$10$UAvOdwl6OeMNAlXbSVwKH.ag86u60RkCUDMTsQHdTnpO7o/msX6SK'

log "Creating auth secret (admin/admin)"
kubectl delete secret pvc-explorer-auth -n pvc-explorer-system --ignore-not-found
kubectl create secret generic pvc-explorer-auth \
  -n pvc-explorer-system \
  --from-literal="admin=${BCRYPT_ADMIN}"

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
