# Overview

**PVC Explorer Operator** is a suite of open-source Kubernetes tools for browsing, managing, and exploring PersistentVolumeClaims — safely, scoped, and on demand.

The suite consists of three projects that work together:

## Projects

### pvc-explorer (Operator)

A Kubernetes-native operator (Kubebuilder v4) that manages ephemeral agent pods for browsing PVCs. It validates PVCs, computes safe mount strategies, enforces read-only fallback when workloads are active, and scales agents to zero when idle.

- [Getting Started](/guide/getting-started) — install and first run
- [Architecture](/architecture) — runtime internals
- [API Reference](/api/rest) — REST, WebSocket, and CRD APIs

### pvc-explorer-agent

A lightweight HTTP file-browser that mounts a PVC and exposes its contents over a REST API. It has an embedded Vue UI and automatically detects conflicts — if another workload is using the same PVC, write endpoints are disabled.

This is the container image deployed by the operator. You don't run it directly — the operator creates agent pods from this image.

- [Agent Overview](/pvc-explorer-agent/overview)

### kubectl-pvc-explorer

A kubectl plugin for exploring PVCs from the command line. List files, view contents, upload/download, execute commands, and even FUSE-mount PVCs — all without the web UI.

- [CLI Overview](/kubectl-pvc-explorer/overview)

## How They Fit Together

<ArchitectureDiagram />

## Getting Started

1. **Install the operator** in your cluster: [Install Guide](/install)
2. **Define scopes** to control which PVCs are accessible
3. **Use the web UI** or **install the CLI plugin** to browse PVCs
4. **Let the operator handle lifecycle** — agents scale to zero when idle
