# ADR 004: Consumer diff uses a union of old and new keys

## Status

Accepted

## Context

`consumer.Index.Sync` needs to emit `consumer.attached` and `consumer.detached` events when the set of pods consuming a PVC changes. It receives the new state from a pod rescan and must compare it against the previous state.

A naive approach — iterate the new state and check for additions, then iterate the old state and check for removals — fails when a PVC key is brand-new (not present in the old map). A range over the old map for the removal check would never visit the new key, so its pods would never be diffed.

## Decision

Build the union of all `(namespace, pvcName)` keys across both the old and new states, then diff each key:

```go
keys := map[string]bool{}
for k := range snapshot { keys[k] = true }
for k := range next     { keys[k] = true }

for k := range keys {
    diffAdded(snapshot[k], next[k])    // pods in next but not in snapshot
    diffRemoved(snapshot[k], next[k])  // pods in snapshot but not in next
}
```

The snapshot is taken inside the write lock before the store is replaced. Events are published after the lock is released.

## Consequences

- Correctly handles new PVC keys (first pod attach) and deleted PVC keys (last pod detach).
- Events are published outside the lock, so the broadcaster's fan-out does not hold the consumer index mutex.
- The diff operates on slices of `ConsumerInfo` structs, each identified by `podName`. A pod appearing in both old and new state under the same PVC key is not reported.
