# ADR 002: Idle countdown is server-pushed via WebSocket

## Status

Accepted

## Context

Agents must scale to zero after `idleTimeout` minutes without a user heartbeat. The UI needs to show a countdown and warn the user before expiry.

Two approaches were considered:

**A. Client-side timer.** Each browser tab tracks elapsed time locally. Sends a heartbeat to reset. Warns and redirects on its own countdown expiry.

**B. Server-side timer with push.** The server tracks the deadline. A background goroutine broadcasts the remaining time every 5 seconds. The browser only renders what the server sends.

## Decision

Server-side timer (option B). The idle watcher goroutine (`internal/scaler.RunIdleWatcher`) is the single source of truth. It:

- Broadcasts `idle.tick` every 5 seconds with the remaining seconds.
- Broadcasts `idle.warning` once when remaining seconds drops to or below 60.
- Calls `SleepAgent` and broadcasts `idle.expired` when remaining reaches zero.

## Consequences

- All tabs viewing the same explorer see the same countdown — no drift between tabs.
- The server enforces the deadline even if the browser tab crashes or the JS timer is inaccurate.
- `idle.warning` fires at most once per wake session (gated by `warnState` map, reset on wake/heartbeat), so a reconnecting tab does not get a duplicate warning.
- The server does 5-second work proportional to the number of Running explorers. This is acceptable at the expected scale (tens of explorers, not thousands).
- Clients do not need to track absolute time — they receive `remainingSeconds` directly and render it.
