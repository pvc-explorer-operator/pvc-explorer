# Overview

PVC Explorer is a Kubernetes operator that gives controlled, on-demand access to PersistentVolumeClaims through an embedded API server and web UI.

## What it does

- Discovers PVCs in the cluster through `PVCExplorerScope` CRDs
- Creates and manages `PVCExplorer` resources for each discovered PVC
- Keeps explorer workloads scaled to zero when idle
- Wakes explorers on demand by deploying a [`pvc-explorer-agent`](https://github.com/pvc-explorer-operator/pvc-explorer-agent) pod that mounts and exposes the PVC
- Proxies file-browser traffic from the UI to the running agent pod
- Exposes REST and WebSocket APIs used by the UI

> **Note:** The operator manages the agent lifecycle but does not contain the agent implementation. See the [`pvc-explorer-agent`](https://github.com/pvc-explorer-operator/pvc-explorer-agent) repository for the file-browser implementation.

## Start here

- Install and first run: [Install](/install)
- Fast local workflow: [Run Local](/guide/local-run)
- Runtime internals: [Architecture](/architecture)
- API contract: [API](/api/)
- Frontend product docs: [UI](/ui/)

## Source references

- Main project entry: https://github.com/pvc-explorer-operator/pvc-explorer
- Existing docs folder: https://github.com/pvc-explorer-operator/pvc-explorer/tree/main/docs
