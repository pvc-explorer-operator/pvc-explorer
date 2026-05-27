# PVCExplorerScope

PVCExplorerScope is a cluster-scoped CRD that registers namespaces and applies defaults used to create PVCExplorer resources.

## API Specification

| Field | Value |
| --- | --- |
| API Version | pvcexplorer.io/v1alpha1 |
| Kind | PVCExplorerScope |
| Scope | Cluster |
| Short Name | pvcs |
| Status Subresource | Enabled |

## Example

```yaml
apiVersion: pvcexplorer.io/v1alpha1
kind: PVCExplorerScope
metadata:
  name: production
spec:
  namespaces:
    names:
      - prod
      - staging
  discovery:
    mode: Auto
    excludePVCs:
      - tmp-*
  deletionPolicy: Cleanup
  defaults:
    mode: ScaledToZero
    forceRW: true
    scaling:
      idleTimeout: 10m
      startupTimeout: 60s
    mountStrategy:
      allowNodeAffinity: true
      fallbackOnConflict: Pending
```

## Fields Reference

### Spec

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| namespaces.names | array of string | No | Static namespace list |
| namespaces.labelSelector | object | No | Label selector for dynamic namespace registration |
| discovery.mode | string | No | PVC discovery mode: Auto or Explicit |
| discovery.pvcNames | array of string | No | Used when mode is Explicit |
| discovery.excludePVCs | array of string | No | Glob patterns excluded from discovery |
| deletionPolicy | string | No | Cleanup or Orphan behavior when scope is deleted |
| defaults.mode | string | No | Default PVCExplorer mode |
| defaults.image | string | No | Default agent image |
| defaults.forceRW | boolean | No | Default read-write preference when safe |
| defaults.scaling.idleTimeout | string | No | Default idle timeout |
| defaults.scaling.startupTimeout | string | No | Default startup timeout |
| defaults.mountStrategy.allowNodeAffinity | boolean | No | Default node affinity permission |
| defaults.mountStrategy.fallbackOnConflict | string | No | Default conflict policy, currently Pending |
| defaults.resources | object | No | Default pod resource requests and limits |

### Status

| Field | Type | Description |
| --- | --- | --- |
| namespaceCount | integer | Number of namespaces currently registered |
| explorerCount | integer | Number of owned PVCExplorer resources |
| observedGeneration | integer | Generation observed by controller |
| conditions | array | Standard Kubernetes conditions |

## Source of truth

- https://github.com/pvc-explorer-operator/pvc-explorer/blob/main/api/v1alpha1/pvcexplorerscope_types.go
