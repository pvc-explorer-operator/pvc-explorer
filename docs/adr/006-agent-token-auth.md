# ADR 006: Agent token authentication for operator→agent isolation

## Status

Accepted

## Context

The operator deploys `pvc-exporter-agent` as a sidecar Deployment per PVCExplorer. The agent exposes an HTTP API (file browser, upload, edit, delete) on a ClusterIP Service. Previously the agent had **no authentication** — any pod on the cluster network could reach the agent API and manipulate filesystem contents.

The agent is also designed to run standalone (outside the operator), where no authentication is appropriate for local development.

We needed:

- **100% isolation**: only the operator can talk to its managed agent instances.
- **Backward compatibility**: standalone mode remains unauthenticated.
- **No permanent secrets**: the auth token should only exist while the agent pod exists.
- **No external dependencies**: the solution must work within the K8s cluster without external auth services.

## Decision

### Two-mode auth in the agent

The agent reads an `AUTH_TOKEN` environment variable at startup:

- **Token set (operator deployment)**: every HTTP request must carry `Authorization: Bearer <token>`. Validation uses `crypto/subtle.ConstantTimeCompare` to prevent timing side-channels.
- **Token empty (standalone)**: the middleware is a no-op pass-through — all requests proceed without auth.

A single `WithAuth(next http.Handler, token string)` middleware wraps the entire `http.ServeMux`, protecting all routes including the UI, API, and health endpoints.

### Operator-managed token lifecycle

The reconciler creates and deletes a K8s Secret (`{explorer-name}-agent-token`) in lockstep with the agent Deployment:

| Agent state | Secret exists? |
|---|---|
| Scaled to zero (no pod) | No — deleted by reconciler |
| Deployment mode (pod running) | Yes — created with random 32-byte hex token |

Lifecycle rules:

- **Wake** (0→1 replicas): if Secret does not exist, generate a fresh token via `crypto/rand` and create it. The Secret carries an ownerReference to the PVCExplorer for cleanup on resource deletion.
- **Sleep** (1→0 replicas): delete the Secret. No pod means no token needed.
- **Reconcile (steady state)**: if Secret and Deployment already exist, no-op.
- **PVCExplorer deletion**: ownerReference cascades to delete the Secret.

### Operator proxy injects the token

The `proxyDispatch` handler in `internal/api/rest.go` forwards requests from the controller's REST API to the agent. Before proxying, it:

1. Reads the token from the `{name}-agent-token` Secret in the explorer's namespace.
2. Sets `Authorization: Bearer <token>` on the cloned request.
3. Forwards via `httputil.ReverseProxy`.

### RBAC

The operator ServiceAccount requires `secrets: get/list/watch/create/update/patch/delete` in managed namespaces to manage the token Secret.

### Key files

| File | Change |
|---|---|
| `pvc-exporter-agent/agent/auth.go` | `WithAuth` middleware (constant-time Bearer check) |
| `pvc-exporter-agent/cmd/agent/main.go` | Reads `AUTH_TOKEN` env, wraps mux |
| `pvc-explorer/internal/controller/agent_reconciler.go` | `reconcileAgentTokenSecret` — create/delete lifecycle |
| `pvc-explorer/internal/api/rest.go` | `readAgentToken` + proxy header injection |
| `pvc-explorer/config/rbac/role.yaml` | Secrets RBAC (regenerated via `make manifests`) |

## Consequences

### Positive

- **Full isolation**: other pods cannot authenticate to the agent. The token is a cryptographically random 64-char hex string, read from a Secret only the operator can access.
- **Zero standing secrets**: the Secret exists only while the agent pod is running. Scaled-to-zero explorers have no Secret anywhere in the cluster.
- **Token rotation on every wake cycle**: sleeping and re-waking generates a fresh token.
- **Backward compatible**: standalone agents without `AUTH_TOKEN` work identically to before.
- **Auto-cleanup**: ownerReference deletes the Secret when the PVCExplorer is removed. No orphan secrets.
- **No new dependencies**: uses only stdlib (`crypto/rand`, `crypto/subtle`, `encoding/hex`) and the K8s Secret API.

### Negative

- **Secret visibility**: the token lives in a K8s Secret. Any principal with `get secret` RBAC in the namespace can read it. In practice, the default `default` ServiceAccount has no such permissions, so this is limited to cluster administrators.
- **Secret per active agent**: each running agent has one Opaque Secret. This is negligible — Secrets have trivial overhead.
- **No token rotation without recycle**: the token is stable for the agent's entire running lifetime. It only changes on scale-down + scale-up. Extending to periodic rotation would require a controller or webhook.
