# pvc-explorer Helm Chart

This chart deploys the pvc-explorer controller and UI as a single deployment.

## Usage

```bash
# Add repo (if published as classic Helm repo)
# helm repo add pvc-explorer https://pvc-explorer-operator.github.io/pvc-explorer/charts
# helm repo update

# OR install directly from a local checkout
helm install pvc-explorer ./helm/pvc-explorer \
  --namespace pvc-explorer-system --create-namespace \
  --set image.tag=<tag>

# OR install from OCI registry (if published)
# helm install pvc-explorer oci://ghcr.io/pvc-explorer-operator/pvc-explorer-chart --version <tag>
```

## Values

- `image.repository` (default: `ghcr.io/pvc-explorer-operator/pvc-explorer`)
- `image.tag` (default: `latest`)
- `replicaCount` (default: 1)
- `resources` (default: `{}`)
- `serviceAccount.create` (default: true)
- `serviceAccount.name` (default: auto)
