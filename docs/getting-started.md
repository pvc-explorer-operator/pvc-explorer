# Getting Started

This guide covers local setup and the first cluster run for pvc-explorer.

## Prerequisites

- Kubernetes v1.24+
- Go 1.24+ for local development
- Docker for image builds
- [kind](https://kind.sigs.k8s.io/) for the recommended local cluster workflow

## Option A — run inside the cluster (recommended)

This is the normal workflow: the controller runs as a pod inside kind.

```bash
# 1. Create the cluster and install CRDs
kind create cluster --config kind/cluster.yaml
make install

# 2. Build and load the image into kind, then deploy
make docker-build IMG=pvc-explorer:dev
kind load docker-image pvc-explorer:dev --name pvc-explorer
make deploy IMG=pvc-explorer:dev

# 3. Apply sample resources
kubectl apply -k config/samples/

# 4. Open the dashboard
kubectl port-forward -n pvc-explorer-system deployment/pvc-explorer-controller-manager 8080:8080
# http://localhost:8080  —  default credentials: admin / admin
```

For subsequent code changes use `kind/rebuild.sh` instead of repeating steps 2–4:

```bash
bash kind/rebuild.sh controller   # rebuild + roll out controller only
bash kind/rebuild.sh              # rebuild controller and agent
```

## Option B — run out-of-cluster (fast dev loop)

The controller runs on your host and talks to the cluster via your kubeconfig.
No image build required — changes are picked up immediately.

```bash
kind create cluster --config kind/cluster.yaml
make install
make run   # blocks — controller logs stream here
```

In a **new terminal**:

```bash
kubectl apply -k config/samples/
# UI is served directly by the controller — no port-forward needed
open http://localhost:8080   # default credentials: admin / admin
```

## kind/ helper scripts

| Script | Purpose |
|---|---|
| `kind/setup.sh` | Creates the cluster, installs CRDs, loads demo namespaces, PVCs and scopes |
| `kind/rebuild.sh` | Rebuilds the controller (and optionally the agent) image and rolls it out |
| `kind/teardown.sh` | Deletes the kind cluster |

## Install in-cluster (production / CI)

kubectl apply -f https://raw.githubusercontent.com/pvc-explorer-operator/pvc-explorer/<tag>/dist/install.yaml
### Option 1: Deploy from public container image

```bash
export IMG=ghcr.io/pvc-explorer-operator/pvc-explorer:<tag>
make deploy IMG=$IMG
```

### Option 2: Helm chart (local or OCI)

```bash
# From local checkout
helm install pvc-explorer ./helm/pvc-explorer \
  --namespace pvc-explorer-system --create-namespace \
  --set image.tag=<tag>

# (Planned) From OCI registry
# helm install pvc-explorer oci://ghcr.io/pvc-explorer-operator/pvc-explorer-chart --version <tag>
```

See [helm/pvc-explorer/README.md](../helm/pvc-explorer/README.md) for chart details.

## Authentication

Credentials are stored in a Kubernetes Secret. See the README for the sample Secret and ConfigMap.

## Registering PVCs

The recommended path is via `PVCExplorerScope`, which auto-discovers Bound PVCs in the configured namespaces and creates a `PVCExplorer` resource for each one.

## Theming

Runtime theming is configured via ConfigMap. See the README for the sample theme payload.
