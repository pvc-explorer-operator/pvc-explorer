# Architecture

PVC Explorer Operator follows a three-zone architecture that separates user access, control logic, and data access.

<ArchitectureDiagram />

## User Space

The user-facing layer consists of two interfaces:

- **Web Dashboard** — Vue-based UI served through the operator's REST API. Provides browser-based file browsing, editing, and management.
- **kubectl-pvc-explorer** — CLI plugin using your existing kubeconfig for authentication. No separate login required.

Both interfaces communicate with the operator's API server through authenticated channels.

## Control Plane

The operator (built with Kubebuilder v4) runs as a Kubernetes deployment and manages the lifecycle of agent pods:

- Validates PVCs and consumer workloads before mounting decisions
- Creates ephemeral agent pods from the `pvc-explorer-agent` image
- Enforces RBAC through namespace-scoped `PVCExplorerScope` resources
- Proxies API requests from the web UI to agent pods
- Scales agent pods to zero when idle (scale-to-zero)

## Data Plane

When a user browses a PVC, the operator creates an agent pod that:

- Mounts the target PVC directly
- Serves a REST file-browser API with read/write and read-only modes
- Watches the Kubernetes API for other pods using the same PVC — if a conflict is detected, write endpoints are disabled (read-only fallback)
- Exposes a `/healthz` probe for readiness checks

## Auth Flow

Authentication flows through the operator's controller proxy:

1. User authenticates via Basic Auth to the operator API
2. Operator validates credentials against the configured auth backend
3. Proxy forwards authenticated requests to the agent pod
4. Agent pod serves file content scoped to the PVC mount

## Lifecycle

Agent pods follow a scale-to-zero lifecycle:

- **Wake** — operator creates an agent pod for the requested PVC
- **Active** — pod serves file content; operator proxies requests
- **Idle timeout** — after a configurable period of inactivity, operator scales the pod to zero
- **Manual sleep** — user can explicitly scale down via the CLI or UI

This architecture ensures PVCs are only mounted when actively accessed, reducing resource usage and attack surface.
