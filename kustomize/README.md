# Kustomize Layout

This directory contains kustomize manifests for local and development-focused deployment workflows.

## Structure

- `base/`: shared resources and defaults
- `overlays/dev/`: local developer overlay (for example `pvc-explorer:dev` image)

## Typical usage

```bash
kubectl apply -k kustomize/overlays/dev/
```

For kubebuilder-managed install/deploy manifests, see `../config/` and project `Makefile` targets.