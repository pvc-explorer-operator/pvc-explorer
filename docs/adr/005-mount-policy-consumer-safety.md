# ADR 005: Mount policy enforces read-only when consumers are present

## Status

Accepted

## Context

A PVC with `accessModes: [ReadWriteOnce]` can only be mounted by pods on the same node. If an application pod already has the PVC mounted read-write, attaching an agent on a different node — or mounting it read-write alongside the existing mount — would cause the Deployment to hang in `Pending` (multi-attach error) or risk data corruption.

## Decision

`mountPolicy` in `agent_reconciler.go` inspects the consumer index before building the Deployment spec:

1. If any consumer pod is Running or Pending, the agent mounts the PVC `readOnly: true`.
2. If `spec.forceRW: true` and no consumers are present, the agent mounts read-write.
3. `targetNode` is set to the node of the existing consumer pod when `accessModes` is `ReadWriteOnce`, so the agent pod is scheduled on the same node and can mount alongside the application without a multi-attach error.

When the last consumer detaches, the `consumer.detached` event triggers a reconcile. If `forceRW: true`, the reconciler patches the Deployment to switch to read-write.

## Consequences

- No data corruption risk from simultaneous read-write mounts.
- The agent remains available (read-only) while the PVC is in use, which is the most useful behavior for debugging running workloads.
- `forceRW: true` is opt-in. The default is read-only when consumers are present.
- Node affinity pinning is only applied for `ReadWriteOnce` PVCs. `ReadWriteMany` PVCs do not need it.
- If a consumer pod moves nodes (e.g., eviction + reschedule), the reconciler will update `targetNode` on the next pod watch event. There may be a brief window where the agent pod is on the wrong node; it will be rescheduled automatically.
