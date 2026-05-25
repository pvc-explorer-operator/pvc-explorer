# ADR 003: agent.ready event is emitted on phase transition in the reconciler

## Status

Accepted

## Context

The WakeUpDialog needs to know when an agent has started successfully so it can navigate the user to the file browser. Two approaches were considered:

**A. HTTP polling.** The dialog polls `GET /api/v1/explorers/{ns}/{name}` every N seconds and checks `status.phase`.

**B. WebSocket event.** The server emits `agent.ready` when the phase transitions to `Running`. The dialog subscribes and navigates on receipt.

For option B, the event could be emitted from:
- The wake REST handler — fires immediately when the wake request is accepted, before the agent is ready.
- The reconciler — fires when the Deployment is actually ready and the phase becomes `Running`.

## Decision

Option B, emitted from the reconciler (`syncStatus` in `agent_reconciler.go`).

```go
prevPhase := explorer.Status.Phase
// ... patch status ...
if prevPhase != ExplorerPhaseRunning && explorer.Status.Phase == ExplorerPhaseRunning {
    broadcaster.Publish("agent.ready", ...)
}
```

## Consequences

- `agent.ready` is accurate: it fires when the pod is Ready, not when the wake request lands.
- Fires exactly once per real transition. The reconciler may be called many times while the deployment is warming up; the phase only crosses the threshold once.
- The WakeUpDialog eliminates its 2-second HTTP polling loop, reducing load and latency.
- A 120-second deadline fallback in the dialog handles the edge case where the event is missed (e.g., reconnect after a long disconnect that exhausted the ring buffer).
