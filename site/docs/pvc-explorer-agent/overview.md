---
title: pvc-explorer-agent
---

# pvc-explorer-agent

A lightweight HTTP file-browser agent that mounts a PersistentVolumeClaim and exposes its contents over a simple REST API.

The agent is the container image that the [pvc-explorer operator](/guide/getting-started) deploys as ephemeral pods. It has an embedded Vue UI for standalone browser-based use and an API that the operator proxies through to the web dashboard.

> By default the agent has no authentication and is only reachable through the controller proxy, which enforces Basic Auth and role checks. Never expose the agent port directly.

## Relationship to the Operator

<div class="rel-chain">

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
  <span class="rel-arrow-label">creates &amp; manages</span>
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

<div class="rel-arrow">
  <svg viewBox="0 0 36 16" fill="none">
    <path d="M2 8h28M28 4l4 4-4 4" stroke="var(--vp-c-brand-1)" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
  </svg>
  <span class="rel-arrow-label">mounts &amp; serves</span>
</div>

<div class="rel-node">

  <svg class="rel-icon rel-icon--amber" viewBox="0 0 44 44" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
    <rect x="7" y="5" width="30" height="34" rx="4" stroke="currentColor"/>
    <path d="M7 13h30" stroke="currentColor"/>
    <rect x="11" y="18" width="8" height="4" rx="1" stroke="currentColor" stroke-width="1.3"/>
    <rect x="11" y="25" width="8" height="4" rx="1" stroke="currentColor" stroke-width="1.3"/>
    <rect x="11" y="32" width="8" height="4" rx="1" stroke="currentColor" stroke-width="1.3"/>
    <path d="M25 20h6M25 27h6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
  </svg>
  <div class="rel-label">PVC</div>
  <div class="rel-sub">Kubernetes storage</div>
</div>

</div>

The operator is responsible for:
- Creating agent pods from the agent container image
- Configuring them with the target PVC name and mount path
- Proxying API requests from the web UI to the agent
- Scaling agents to zero when idle

The agent itself has no awareness of the operator — it simply mounts a PVC and serves an HTTP API.

## Features

- HTTP file browser endpoints for listing, downloading, editing, and uploading files
- Read-only fallback — watches the Kubernetes API for other pods using the same PVC; if a conflict is detected, write endpoints are disabled
- Embedded Vue UI for standalone browser-based use
- `/healthz` probe for readiness checks
- Optional Bearer token auth via `AUTH_TOKEN` environment variable

## Container Images

OCI images are published at `ghcr.io/pvc-explorer-operator/pvc-explorer-agent`:

| Tag         | Description                           |
| ----------- | ------------------------------------- |
| `<version>` | Specific release (e.g. `v0.1.0`)      |
| `dev`       | Mutable development image from `main` |
| `latest`    | Latest stable release                 |

Images are multi-arch (`linux/amd64`, `linux/arm64`), cosign-signed with BuildKit provenance and SBOM artifacts.

## Quick Start

Run the agent locally with a demo dataset:

```shell
make run-agent ROOT=./testdata/demo PVC=demo-pvc
```

The repository includes a small demo dataset in `./testdata/demo` for local UI and API testing.

## Repository

- **Source:** [github.com/pvc-explorer-operator/pvc-explorer-agent](https://github.com/pvc-explorer-operator/pvc-explorer-agent)
- **License:** Apache 2.0
