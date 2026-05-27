# Run Local

This guide covers two local workflows:

- Run the UI locally with Vite
- Run the full stack in a local kind cluster using the helper scripts in `kind/`

## Run the UI locally

### Fast UI development with mock data

Use this workflow for screens, layout, state handling, and UI behavior without a cluster.

```bash
cd ui
npm install
npm run dev
```

In development mode, the UI uses the local mock plugin and auto-authenticates for fast iteration.

### Test the real login flow locally

Use this workflow to test the login page against a live backend instead of the development auth bypass.

Terminal 1:

```bash
cd ui
VITE_DEV_AUTH_BYPASS=false npm run dev
```

Terminal 2:

```bash
make run
```

Then open the UI and navigate to `/login`.

## Run the full stack in kind

The repo includes helper scripts in `kind/` for a repeatable local cluster workflow.

### Prerequisites

Make sure these tools are installed and available on your `PATH`:

- `kind`
- `kubectl`
- `docker`
- `make`

### Create and populate the demo cluster

From the repository root:

```bash
kind/setup.sh
```

What `kind/setup.sh` does:

- creates the `pvc-explorer` kind cluster if it does not already exist
- builds the controller image as `pvc-explorer:dev`
- pulls the agent image from `ghcr.io/pvc-explorer-operator/pvc-explorer-agent:dev` (see [pvc-explorer-agent](https://github.com/pvc-explorer-operator/pvc-explorer-agent))
- loads both images into kind
- regenerates and applies CRDs
- deploys the controller using the dev overlay in `kustomize/overlays/dev`
- creates the `pvc-explorer-auth` secret with `admin / admin`
- applies the demo storage class, PVs, PVCs, namespaces, and example scopes

When setup finishes, the local dashboard is available at `http://localhost:8080`.

Default login:

- username: `admin`
- password: `admin`

Demo namespaces and scopes created by the script:

- `demo`
- `demo-staging`
- `PVCExplorerScope/demo`
- `PVCExplorerScope/demo-by-label`

Useful checks:

```bash
kubectl get pvcexplorerscope,pvcexplorer -A
kubectl logs -n pvc-explorer-system -l control-plane=controller-manager -f
```

### Rebuild after code changes

Use the rebuild helper instead of recreating the whole cluster.

Rebuild controller only:

```bash
kind/rebuild.sh controller
```

Reload agent image only:

```bash
kind/rebuild.sh agent
```

Rebuild both:

```bash
kind/rebuild.sh
```

Notes:

- the controller rebuild loads a new image tag into the cluster and updates the deployment
- the agent rebuild path reloads the image into kind; existing agent pods pick it up on the next wake cycle

### Override the agent image

If the default GHCR image is not accessible, provide a different agent image explicitly:

```bash
AGENT_IMG=<your-agent-image> kind/setup.sh
```

or:

```bash
AGENT_IMG=<your-agent-image> kind/rebuild.sh agent
```

### Tear the cluster down

To remove the kind cluster and its local data:

```bash
kind/teardown.sh
```

The script prompts for confirmation before deleting the `pvc-explorer` cluster.
