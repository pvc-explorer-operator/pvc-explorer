# ADR 001: Broadcaster uses string event type to avoid import cycles

## Status

Accepted

## Context

Three packages need to publish WebSocket events:

- `internal/api` — REST mutation responses
- `internal/scaler` — idle watcher goroutine
- `internal/controller` — agent phase transitions

`internal/api` owns the concrete `Broadcaster` type and the `EventType` type alias. If `internal/scaler` or `internal/controller` imported `internal/api` to call `Publish(EventType, any)`, they would create an import cycle (`api` imports `scaler` for `WakeAgent`/`SleepAgent`; `api` imports `controller` for the reconcilers).

## Decision

`Broadcaster.Publish` accepts `string` instead of `EventType`:

```go
func (b *Broadcaster) Publish(eventType string, payload any) error
```

Inside `ws.go`, the string is cast to `EventType` when building the frame.

Packages outside `internal/api` define a minimal local interface with the same signature:

```go
type Broadcaster interface {
    Publish(eventType string, payload any) error
}
```

`*api.Broadcaster` satisfies these interfaces. `cmd/main.go` constructs one instance and passes it to all consumers.

## Consequences

- No import cycles.
- Event type constants in `ws_types.go` are the authoritative list. Callers outside `internal/api` use string literals. A typo produces a silent no-op rather than a compile error. Mitigate by keeping the constants in `ws_types.go` and referencing them from a shared constants file if the project grows.
- Adding a new event does not require changing the interface in every package.
