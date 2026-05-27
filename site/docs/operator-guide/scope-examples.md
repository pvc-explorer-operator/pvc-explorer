# Scope Examples

This page shows two common patterns for exposing PVCs through PVC Explorer:

1. **Explicit namespace** — register a namespace by name in the scope
2. **Label-selected namespace** — let the scope discover namespaces dynamically by label

---

## Pattern 1: Explicit namespace

Use this pattern when you know the exact namespace and want precise control over which namespaces a scope covers.

### Create the namespace and PVC

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: my-app
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: my-app-data
  namespace: my-app
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
```

```bash
kubectl apply -f namespace-and-pvc.yaml
```

### Create the scope

```yaml
apiVersion: pvcexplorer.io/v1alpha1
kind: PVCExplorerScope
metadata:
  name: my-app-scope
spec:
  namespaces:
    names:
      - my-app
  discovery:
    mode: Auto
  deletionPolicy: Cleanup
  defaults:
    mode: ScaledToZero
    image: ghcr.io/pvc-explorer-operator/pvc-explorer-agent:latest
    scaling:
      idleTimeout: "10m"
      startupTimeout: "60s"
    mountStrategy:
      allowNodeAffinity: true
      fallbackOnConflict: Pending
    resources:
      requests:
        cpu: "50m"
        memory: "64Mi"
      limits:
        cpu: "200m"
        memory: "256Mi"
```

```bash
kubectl apply -f scope.yaml
```

### What happens

The scope reconciler finds `my-app` in `spec.namespaces.names` and creates a `PVCExplorer` resource for `my-app-data`. The explorer starts in `ScaledToZero` phase. Use the UI or the REST API to wake it on demand.

```bash
kubectl get pvcexplorer -n my-app
```

---

## Pattern 2: Label-selected namespace

Use this pattern when namespaces are created dynamically or you want to opt multiple namespaces in by adding a single label, without updating the scope each time.

### Create the namespace with the opt-in label and a PVC

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: my-team
  labels:
    pvc-explorer: enabled
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: my-team-data
  namespace: my-team
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 2Gi
```

```bash
kubectl apply -f namespace-and-pvc.yaml
```

The label `pvc-explorer: enabled` on the namespace is what the scope watches. Any namespace carrying this label is automatically enrolled when the scope reconciler runs.

### Create the scope

```yaml
apiVersion: pvcexplorer.io/v1alpha1
kind: PVCExplorerScope
metadata:
  name: team-scope
spec:
  namespaces:
    labelSelector:
      matchLabels:
        pvc-explorer: enabled
  discovery:
    mode: Auto
    excludePVCs:
      - "tmp-*"
  deletionPolicy: Cleanup
  defaults:
    mode: ScaledToZero
    image: ghcr.io/pvc-explorer-operator/pvc-explorer-agent:latest
    scaling:
      idleTimeout: "10m"
      startupTimeout: "60s"
    mountStrategy:
      allowNodeAffinity: true
      fallbackOnConflict: Pending
    resources:
      requests:
        cpu: "50m"
        memory: "64Mi"
      limits:
        cpu: "200m"
        memory: "256Mi"
```

```bash
kubectl apply -f scope.yaml
```

### What happens

The scope reconciler watches all namespaces. When it finds `my-team` (labelled `pvc-explorer: enabled`), it creates a `PVCExplorer` for `my-team-data`. Add the same label to any new namespace and it is enrolled automatically — no scope update required.

```bash
# Opt a new namespace in at any time
kubectl label namespace another-team pvc-explorer=enabled

# Verify explorers were created
kubectl get pvcexplorer -A
```

---

## Verify the full picture

```bash
kubectl get pvcexplorerscope
kubectl get pvcexplorer -A
```

Expected output example:

```
NAME           AGE
my-app-scope   2m
team-scope     1m

NAMESPACE   NAME           PHASE          AGE
my-app      my-app-data    ScaledToZero   2m
my-team     my-team-data   ScaledToZero   1m
```

---

## Reference

- [PVCExplorerScope CRD](/api/crds/pvcexplorerscope)
- [PVCExplorer CRD](/api/crds/pvcexplorer)
- [`pvc-explorer-agent`](https://github.com/pvc-explorer-operator/pvc-explorer-agent) — the file-browser agent deployed by the operator
