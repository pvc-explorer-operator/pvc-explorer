# PVCExplorer

PVCExplorer is a namespaced CRD that manages the lifecycle of one [`pvc-explorer-agent`](https://github.com/pvc-explorer-operator/pvc-explorer-agent) pod for one PersistentVolumeClaim. The agent is a standalone file-browser implementation deployed and proxied by the operator.

## API Specification

| Field              | Value                   |
| ------------------ | ----------------------- |
| API Version        | pvcexplorer.io/v1alpha1 |
| Kind               | PVCExplorer             |
| Scope              | Namespaced              |
| Short Name         | pvcexp                  |
| Status Subresource | Enabled                 |

## Example

```yaml
apiVersion: pvcexplorer.io/v1alpha1
kind: PVCExplorer
metadata:
  name: app-data
  namespace: default
spec:
  pvcName: app-data
  mode: ScaledToZero
  forceRW: true
  scaling:
    idleTimeout: 10m
    startupTimeout: 60s
  mountStrategy:
    autoDetect: true
    allowNodeAffinity: true
    fallbackOnConflict: Pending
```

## Fields Reference

### Spec

| Field                            | Type              | Required | Description                                                            |
| -------------------------------- | ----------------- | -------- | ---------------------------------------------------------------------- |
| explorerLabels                   | map[string]string | No       | Extra labels for Deployment and pods                                   |
| forceRW                          | boolean           | No       | Use read-write when possible; can be forced read-only by safety checks |
| image                            | string            | No       | Agent container image override                                         |
| mode                             | string            | No       | Agent mode: ScaledToZero or Deployment                                 |
| mountStrategy.allowNodeAffinity  | boolean           | No       | Allow node affinity for RWO mount scenarios                            |
| mountStrategy.autoDetect         | boolean           | No       | Detect active consumers and enforce read-only safety                   |
| mountStrategy.fallbackOnConflict | string            | No       | Current supported value: Pending                                       |
| port                             | integer           | No       | Agent HTTP port, default 8081                                          |
| pvcName                          | string            | Yes      | Target PVC name in the same namespace                                  |
| resources                        | object            | No       | Pod resource requests and limits                                       |
| scaling.idleTimeout              | string            | No       | Inactivity duration before scale-down                                  |
| scaling.provider                 | string            | No       | auto, native, or knative                                               |
| scaling.startupTimeout           | string            | No       | Max wait for agent readiness after wake                                |
| subPath                          | string            | No       | Optional sub-directory of the PVC                                      |

### Status

| Field                 | Type     | Description                                                     |
| --------------------- | -------- | --------------------------------------------------------------- |
| agent.cluster         | string   | Cluster identity reported by agent                              |
| agent.pvcWatchEnabled | boolean  | Agent PVC watch state                                           |
| agent.version         | string   | Agent version string                                            |
| agentEndpoint         | string   | In-cluster endpoint for the agent service                       |
| conditions            | array    | Standard Kubernetes conditions                                  |
| lastHealthCheck       | datetime | Last successful health probe timestamp                          |
| mode                  | string   | Observed mode                                                   |
| mount.accessMode      | string   | PVC access mode                                                 |
| mount.consumers       | array    | Current consumer pods for the PVC                               |
| mount.forceRWDeferred | boolean  | Indicates delayed return to read-write when consumers detach    |
| mount.readOnly        | boolean  | Current read-only state                                         |
| mount.strategy        | string   | Selected mount strategy                                         |
| mount.targetNode      | string   | Target node when node affinity is used                          |
| observedGeneration    | integer  | Generation observed by controller                               |
| phase                 | string   | Lifecycle phase: ScaledToZero, Waking, Running, Failed, Pending |

## Source of truth

- https://github.com/pvc-explorer-operator/pvc-explorer/blob/main/api/v1alpha1/pvcexplorer_types.go
