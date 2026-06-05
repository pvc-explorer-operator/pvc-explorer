---
title: kubectl-pvc-explorer
---

# kubectl-pvc-explorer

A kubectl plugin for browsing and managing PVC contents via the [pvc-explorer operator](/guide/getting-started).

Use your existing kubeconfig for authentication — no extra login or credentials required.

```shell
kubectl pvc list
kubectl pvc ls my-pvc -n my-namespace /data/logs
kubectl pvc cp my-pvc:/remote/file.txt ./local-file.txt
```

## Relationship to the Operator

The CLI plugin communicates with the PVC Explorer operator installed in your cluster:

<div class="rel-chain">

<div class="rel-node">
  <svg class="rel-icon" viewBox="0 0 40 40" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
    <rect x="4" y="10" width="32" height="22" rx="3" stroke="currentColor"/>
    <path d="M4 16h32" stroke="currentColor"/>
    <circle cx="8" cy="13" r="1" fill="currentColor" stroke="none"/>
    <circle cx="11" cy="13" r="1" fill="currentColor" stroke="none"/>
    <circle cx="14" cy="13" r="1" fill="currentColor" stroke="none"/>
    <path d="M14 22l-4 4 4 4" stroke="currentColor" stroke-width="1.5"/>
    <path d="M22 22l4 4-4 4" stroke="currentColor" stroke-width="1.5"/>
  </svg>
  <div class="rel-label">kubectl pvc</div>
  <div class="rel-sub">CLI plugin</div>
</div>

<div class="rel-arrow">
  <svg viewBox="0 0 36 16" fill="none">
    <path d="M2 8h28M28 4l4 4-4 4" stroke="var(--vp-c-brand-1)" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
  </svg>
  <span class="rel-arrow-label">kubeconfig</span>
</div>

<div class="rel-node">
  <svg class="rel-icon rel-icon--blue" viewBox="0 0 44 44" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
    <circle cx="22" cy="22" r="16" stroke="currentColor"/>
    <circle cx="22" cy="22" r="8" stroke="currentColor"/>
    <path d="M22 2v6M22 36v6M2 22h6M36 22h6" stroke="currentColor" stroke-width="1.5"/>
    <path d="M8.5 8.5l4.2 4.2M31.3 31.3l4.2 4.2M31.3 8.5l-4.2 4.2M8.5 31.3l4.2-4.2" stroke="currentColor" stroke-width="1.3"/>
    <circle cx="22" cy="22" r="3" fill="currentColor" opacity="0.3" stroke="none"/>
  </svg>
  <div class="rel-label">pvc-explorer</div>
  <div class="rel-sub">Operator</div>
</div>

<div class="rel-arrow">
  <svg viewBox="0 0 36 16" fill="none">
    <path d="M2 8h28M28 4l4 4-4 4" stroke="var(--vp-c-brand-1)" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
  </svg>
  <span class="rel-arrow-label">proxies</span>
</div>

<div class="rel-node">
  <svg class="rel-icon rel-icon--cyan" viewBox="0 0 44 44" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
    <rect x="6" y="8" width="32" height="28" rx="4" stroke="currentColor"/>
    <path d="M6 16h32" stroke="currentColor"/>
    <path d="M16 8v28M28 8v28" stroke="currentColor" stroke-width="1.3" opacity="0.4"/>
    <rect x="10" y="20" width="4" height="3" rx="1" fill="currentColor" opacity="0.4" stroke="none"/>
    <rect x="10" y="26" width="4" height="3" rx="1" fill="currentColor" opacity="0.4" stroke="none"/>
    <rect x="10" y="32" width="4" height="2" rx="1" fill="currentColor" opacity="0.4" stroke="none"/>
    <path d="M30 24h4M30 28h4" stroke="currentColor" stroke-width="1.5"/>
  </svg>
  <div class="rel-label">pvc-explorer-agent</div>
  <div class="rel-sub">Agent pod</div>
</div>

</div>

The plugin reads the bearer token from the `*-agent-token` Secret automatically, so there is no separate login step.

## Installation

```shell
git clone https://github.com/pvc-explorer-operator/kubectl-pvc-explorer.git
cd kubectl-pvc-explorer
make install
```

Verify:

```shell
kubectl plugin list
kubectl pvc --version
```

### Prerequisites

- Go 1.26+
- kubectl with a valid kubeconfig pointing to a cluster with PVC Explorer installed
- macOS: [macFUSE](https://github.com/macfuse/macfuse/wiki/Getting-Started) (only for the `mount` command)

## Commands

| Command                           | Description                                                 |
| --------------------------------- | ----------------------------------------------------------- |
| `kubectl pvc cat <pvc> <path>`    | Display file contents                                       |
| `kubectl pvc cp <src> <dest>`     | Upload/download files (supports glob patterns)              |
| `kubectl pvc exec <pvc> -- <cmd>` | Execute commands inside the agent pod                       |
| `kubectl pvc list`                | List all PVCExplorer resources                              |
| `kubectl pvc ls <pvc> [path]`     | List files on a PVC                                         |
| `kubectl pvc mount <pvc> <dir>`   | FUSE mount (experimental)                                   |
| `kubectl pvc sleep <pvc>`         | Scale agent back to 0                                       |
| `kubectl pvc wake <pvc>`          | Scale agent from 0 to 1 (`--wait 2m` to wait for readiness) |

## Quick Start

```shell
# Wake the agent
kubectl pvc wake my-pvc -n my-namespace

# Browse files
kubectl pvc ls my-pvc -n my-namespace /data

# Download a file
kubectl pvc cp my-pvc:/data/config.yaml ./config.yaml

# Upload a file
kubectl pvc cp ./local-file.txt my-pvc:/remote/dir/

# Execute a command
kubectl pvc exec my-pvc -n my-namespace -- ls -la

# Sleep when done
kubectl pvc sleep my-pvc -n my-namespace
```

## Local Development

The project includes a Kind-based local environment in the [pvc-explorer](https://github.com/pvc-explorer-operator/pvc-explorer) repo:

```shell
cd pvc-explorer
./kind/setup.sh
cd ../kubectl-pvc-explorer
make install
kubectl pvc list
```

## Repository

- **Source:** [github.com/pvc-explorer-operator/kubectl-pvc-explorer](https://github.com/pvc-explorer-operator/kubectl-pvc-explorer)
- **License:** Apache 2.0
