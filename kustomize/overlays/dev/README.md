# Dev overlay for Kustomize

This overlay sets the controller image to the locally built `pvc-explorer:dev` tag for use with kind.

## Usage

1. Build the image and load it into kind:

```bash
docker build -t pvc-explorer:dev .
kind load docker-image pvc-explorer:dev --name pvc-explorer
```

2. Apply the overlay:

```bash
kubectl apply -k kustomize/overlays/dev/
```
