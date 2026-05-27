# Install

This page summarizes the installation paths and links to the source documentation.

## Prerequisites

- Kubernetes cluster (kind recommended for local)
- Go toolchain and Docker
- kubectl and make

## Quick install on kind

```bash
kind create cluster --config kind/cluster.yaml
make docker-build IMG=pvc-explorer:dev
kind load docker-image pvc-explorer:dev --name pvc-explorer
make install && make deploy IMG=pvc-explorer:dev
kubectl apply -k config/samples/
```

Open dashboard:

```bash
kubectl -n pvc-explorer-system port-forward svc/pvc-explorer-controller-manager 8080:8080
```

Then visit `http://localhost:8080`

## Fast dev loop (out-of-cluster manager)

```bash
make install
make run
```

Use [Run Local](/guide/local-run) for complete local loops, including UI mock mode.

## Full source guide

- https://github.com/pvc-explorer-operator/pvc-explorer/blob/main/docs/getting-started.md
