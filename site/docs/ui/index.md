# UI

The PVC Explorer web interface allows users to discover, wake, and interact with PersistentVolumeClaims in a Kubernetes cluster.

Users can browse PVC contents directly in the browser through an integrated file browser, manage explorer lifecycle (wake on demand, scale to zero when idle), and configure scopes to control which PVCs are visible.

> The file-browser functionality is provided by [`pvc-explorer-agent`](https://github.com/pvc-explorer-operator/pvc-explorer-agent), a standalone component deployed and proxied by the operator.

## In this section

| Page                              | What it covers                                 |
| --------------------------------- | ---------------------------------------------- |
| [Components](./components.md)     | Individual UI components and where they appear |
| [Contributing](./contributing.md) | Local development and implementation reference |
| [Flows](./flows.md)               | End-to-end user journeys through the interface |
| [UI Display](./display.md)        | Screens, user actions, and expected behavior   |
