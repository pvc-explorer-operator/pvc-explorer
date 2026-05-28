# CRD Release Automation - Complete Process Documentation

## Overview

This document explains the complete process for automating CRD (Custom Resource Definition) release packaging for the PVC-Explorer Kubernetes operator. It covers both automated releases triggered by git tags and manual local builds.

## Table of Contents

1. [Architecture](#architecture)
2. [Release Workflow](#release-workflow)
3. [Installation Methods](#installation-methods)
4. [Local Development](#local-development)
5. [Troubleshooting](#troubleshooting)
6. [Reference](#reference)

---

## Architecture

### What are CRDs?

Custom Resource Definitions (CRDs) extend Kubernetes by allowing users to define custom resource types. PVC-Explorer provides two CRDs:

- **PVCExplorer** (`pvcexplorers.pvcexplorer.io`): Manages individual PVC browsing sessions with automatic scale-to-zero
- **PVCExplorerScope** (`pvcexplorerscopes.pvcexplorer.io`): Defines scoped access to PVCs across namespaces

### Release Artifacts

Each release produces two YAML manifest files:

```
pvc-explorer-crds.yaml          # CRDs only
pvc-explorer-install.yaml       # Full stack (CRDs + RBAC + Controller)
```

### Build Pipeline

```
Repository Structure
├── config/crd/                 # CRD definitions
│   ├── bases/
│   │   ├── pvcexplorer.io_pvcexplorers.yaml
│   │   └── pvcexplorer.io_pvcexplorerscopes.yaml
│   ├── kustomization.yaml
│   └── kustomizeconfig.yaml
├── config/default/             # Full installation (CRDs + RBAC + deployment)
├── .github/workflows/
│   └── release-crds.yaml       # GitHub Actions workflow
└── scripts/
    └── build-crds.sh           # Local build script
```

---

## Release Workflow

### Automated Release Process

When you push a git tag matching the pattern `v*.*.*`, GitHub Actions automatically:

1. **Triggers the Workflow** (`release-crds.yaml`)
2. **Installs Kustomize** on the runner
3. **Builds Manifests** from `config/crd/` and `config/default/`
4. **Verifies Output** to ensure CRDs are present
5. **Attaches to Release** and creates release notes

### Step-by-Step: Creating a Release

#### 1. Tag and Push

```bash
# Ensure you're on main branch with latest code
git checkout main
git pull origin main

# Create a semantic version tag
git tag v1.0.0

# Push the tag to trigger the workflow
git push origin v1.0.0
```

#### 2. Workflow Executes

The `.github/workflows/release-crds.yaml` workflow:

```yaml
# Triggers on semver tags
on:
  push:
    tags:
      - 'v*.*.*'
```

**Workflow Steps:**

| Step                | Action                           | Output                                    |
| ------------------- | -------------------------------- | ----------------------------------------- |
| 1. Checkout         | Clone repo at tag                | Full source code                          |
| 2. Kustomize Setup  | Install kustomize binary         | `/usr/local/bin/kustomize`                |
| 3. Build CRDs       | `kustomize build config/crd`     | `build/release/pvc-explorer-crds.yaml`    |
| 4. Build Full Stack | `kustomize build config/default` | `build/release/pvc-explorer-install.yaml` |
| 5. Verify           | Check manifest integrity         | Line counts, resource counts              |
| 6. Release          | Upload to GitHub Release         | Release assets, notes                     |

#### 3. Manual Verification

After the workflow completes:

```bash
# Check the release was created
curl -s https://api.github.com/repos/pvc-explorer-operator/pvc-explorer/releases/latest | jq '.tag_name, .assets'

# View the generated YAML
curl -O https://github.com/pvc-explorer-operator/pvc-explorer/releases/download/v1.0.0/pvc-explorer-crds.yaml
curl -O https://github.com/pvc-explorer-operator/pvc-explorer/releases/download/v1.0.0/pvc-explorer-install.yaml
```

---

## Installation Methods

### Installation Option 1: Full Stack (Recommended)

**Best for:** Production deployments where you want everything

```bash
# Replace v1.0.0 with actual release version
kubectl apply -f https://github.com/pvc-explorer-operator/pvc-explorer/releases/download/v1.0.0/pvc-explorer-install.yaml
```

**What gets installed:**

- PVCExplorer CRD
- PVCExplorerScope CRD
- RBAC roles and role bindings
- Controller deployment in `pvc-explorer-system` namespace
- Service account and other supporting resources

**Verification:**

```bash
# Check namespace was created
kubectl get namespace pvc-explorer-system

# Check deployment is running
kubectl get deployment -n pvc-explorer-system
kubectl get pods -n pvc-explorer-system

# Check CRDs registered
kubectl get crd | grep pvcexplorer.io
```

### Installation Option 2: CRDs Only

**Best for:** When you manage controller separately or just want to define CRD resources

```bash
kubectl apply -f https://github.com/pvc-explorer-operator/pvc-explorer/releases/download/v1.0.0/pvc-explorer-crds.yaml
```

**What gets installed:**

- PVCExplorer CRD definition
- PVCExplorerScope CRD definition

**Verification:**

```bash
# Check CRDs are registered
kubectl api-resources | grep pvcexplorer

# Describe a CRD to see its schema
kubectl describe crd pvcexplorers.pvcexplorer.io
```

**Next Steps - Install controller separately:**

```bash
# Build and deploy controller manually
git clone https://github.com/pvc-explorer-operator/pvc-explorer.git
cd pvc-explorer
make docker-build IMG=pvc-explorer:v1.0.0
make deploy IMG=pvc-explorer:v1.0.0
```

### Installation Option 3: From Source with Kustomize

**Best for:** Development, customization, or air-gapped environments

```bash
# Clone repository
git clone https://github.com/pvc-explorer-operator/pvc-explorer.git
cd pvc-explorer

# Checkout specific version
git checkout v1.0.0

# Install CRDs only
kubectl apply -k config/crd

# Or install full stack (with customization)
kubectl apply -k config/default
```

**Customization Example:**

Edit `config/default/kustomization.yaml` to customize:

- Resource limits
- Replica count
- Environment variables
- Labels and annotations

```bash
# Then apply
kubectl apply -k config/default
```

---

## Local Development

### Building CRDs Locally

Use the provided build script for local development:

```bash
# Make script executable
chmod +x scripts/build-crds.sh

# Run build script
./scripts/build-crds.sh
```

**Output:**

```
Building CRD manifests...
Building CRDs only...
✓ Generated: build/pvc-explorer-crds.yaml
Building full install bundle...
✓ Generated: build/pvc-explorer-install.yaml

Summary:
  CRDs: 2 found
  Total resources in CRDs: 2
  Total resources in install: 15

To apply locally:
  kubectl apply -f build/pvc-explorer-crds.yaml
  kubectl apply -f build/pvc-explorer-install.yaml
```

### Manual Build Process

If you prefer to build manually:

```bash
# Install kustomize if not already installed
curl -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh" | bash
sudo mv kustomize /usr/local/bin/

# Create output directory
mkdir -p build

# Build CRDs only
kustomize build config/crd > build/pvc-explorer-crds.yaml

# Build full installation
kustomize build config/default > build/pvc-explorer-install.yaml

# Verify the builds
wc -l build/*.yaml
grep "kind:" build/pvc-explorer-crds.yaml | sort | uniq -c
```

### Local Testing Workflow

```bash
# 1. Make changes to CRD definitions
# Edit: config/crd/bases/pvcexplorer.io_pvcexplorers.yaml

# 2. Rebuild manifests
./scripts/build-crds.sh

# 3. Test with kind cluster
kind create cluster --config kind/cluster.yaml
kubectl apply -f build/pvc-explorer-crds.yaml

# 4. Create test resources
kubectl apply -f - <<EOF
apiVersion: pvcexplorer.io/v1alpha1
kind: PVCExplorer
metadata:
  name: test
spec:
  pvcName: data-pvc
EOF

# 5. Verify
kubectl get pvcexplorer test -o yaml

# 6. Cleanup
kind delete cluster
```

---

## Workflow Files

### GitHub Actions Workflow (`.github/workflows/release-crds.yaml`)

**Triggers:**

- On git tag push matching `v*.*.*` (semver)
- Manual trigger via `workflow_dispatch`

**Key Jobs:**

1. **Checkout** - Gets repository at tag
2. **Kustomize Setup** - Installs build tool
3. **Build Manifests** - Generates YAML files
4. **Verify** - Ensures quality and completeness
5. **Release** - Attaches files and creates notes

**Release Notes Template:**

The workflow auto-generates installation instructions in the release notes:

```markdown
# PVC-Explorer Release vX.Y.Z

## Installation

### Option 1: Full Installation (Recommended)
kubectl apply -f https://github.com/pvc-explorer-operator/pvc-explorer/releases/download/vX.Y.Z/pvc-explorer-install.yaml

### Option 2: CRDs Only
kubectl apply -f https://github.com/pvc-explorer-operator/pvc-explorer/releases/download/vX.Y.Z/pvc-explorer-crds.yaml

### Option 3: Using Kustomize
git clone https://github.com/pvc-explorer-operator/pvc-explorer.git
cd pvc-explorer
git checkout vX.Y.Z
kubectl apply -k config/default
```

### Build Script (`scripts/build-crds.sh`)

Automates local CRD building with error checking:

```bash
#!/bin/bash
# Features:
# - Validates kustomize installation
# - Creates output directory
# - Builds both CRD and full manifests
# - Verifies outputs
# - Provides summary
# - Suggests next steps
```

### Kustomize Configuration (`config/crd/kustomization.yaml`)

Defines what resources to include in the CRD manifest:

```yaml
resources:
- bases/pvcexplorer.io_pvcexplorerscopes.yaml
- bases/pvcexplorer.io_pvcexplorers.yaml

# Patches and configurations for customization
patches: []
```

---

## Troubleshooting

### Problem: Workflow Not Triggering

**Symptoms:** Tagged commit pushed but workflow doesn't run

**Solutions:**

1. **Check tag format** - Must match `v*.*.*`

   ```bash
   # ✓ Correct
   git tag v1.0.0
   git tag v1.2.3-rc1
   
   # ✗ Wrong
   git tag 1.0.0        # Missing 'v'
   git tag release-1.0  # Wrong format
   ```

2. **Verify tag pushed** - Tags must be explicitly pushed

   ```bash
   git push origin v1.0.0
   # NOT: git push (doesn't push tags)
   ```

3. **Check workflow file exists** - Must be in default branch

   ```bash
   git show main:.github/workflows/release-crds.yaml
   ```

4. **Verify permissions** - Need `contents: write` permission

   ```yaml
   permissions:
     contents: write
   ```

### Problem: Kustomize Build Fails

**Symptoms:** Workflow fails at "Build CRD manifests" step

**Solutions:**

1. **Install kustomize locally** and test build

   ```bash
   ./scripts/build-crds.sh
   ```

2. **Check kustomization.yaml syntax**

   ```bash
   kustomize build config/crd > /dev/null
   ```

3. **Verify file paths** - All referenced files must exist

   ```bash
   ls config/crd/bases/pvcexplorer.io_*.yaml
   ```

### Problem: CRDs Not Found After Installation

**Symptoms:** `kubectl get crd` doesn't show PVC Explorer CRDs

**Solutions:**

1. **Verify manifest was applied**

   ```bash
   kubectl apply -f pvc-explorer-crds.yaml --dry-run=client
   ```

2. **Check for errors**

   ```bash
   kubectl apply -f pvc-explorer-crds.yaml
   kubectl get events --all-namespaces | grep -i error
   ```

3. **Inspect CRD object**

   ```bash
   kubectl get crd pvcexplorers.pvcexplorer.io -o yaml
   kubectl describe crd pvcexplorers.pvcexplorer.io
   ```

### Problem: Resource Creation Fails

**Symptoms:** `kubectl apply` fails when creating PVCExplorer resource

**Solutions:**

1. **Verify CRD installed**

   ```bash
   kubectl api-resources | grep -i pvcexplorer
   ```

2. **Check schema validation**

   ```bash
   kubectl explain pvcexplorer.spec
   ```

3. **Validate resource YAML**

   ```bash
   kubectl apply -f resource.yaml --dry-run=client -o yaml
   ```

4. **View validation errors**

   ```bash
   kubectl apply -f resource.yaml -v=8  # Verbose output
   ```

---

## Reference

### File Structure

```
pvc-explorer/
├── .github/
│   └── workflows/
│       └── release-crds.yaml           # GitHub Actions workflow
├── config/
│   ├── crd/
│   │   ├── bases/
│   │   │   ├── pvcexplorer.io_pvcexplorers.yaml      # PVCExplorer CRD
│   │   │   └── pvcexplorer.io_pvcexplorerscopes.yaml # PVCExplorerScope CRD
│   │   ├── kustomization.yaml
│   │   └── kustomizeconfig.yaml
│   ├── default/                        # Full installation config
│   ├── manager/                        # Controller deployment
│   └── rbac/                           # RBAC configuration
├── docs/
│   ├── operations/crd-installation.md  # Installation guide
│   ├── getting-started.md
│   ├── architecture.md
│   └── releases.md
├── scripts/
│   └── build-crds.sh                   # Build script
└── build/                              # Generated manifests (local)
    ├── pvc-explorer-crds.yaml
    └── pvc-explorer-install.yaml
```

### CRD Details

#### PVCExplorer

- **API Group:** `pvcexplorer.io`
- **Kind:** `PVCExplorer`
- **Short Name:** `pvcexp`
- **Scope:** Namespaced
- **Status Subresource:** Yes

**Key Spec Fields:**

- `pvcName` (required) - Name of PVC to explore
- `mode` - Scale-to-zero or always running
- `port` - HTTP port (default: 8081)
- `resources` - CPU/memory limits
- `scaling` - Idle timeout and provider

#### PVCExplorerScope

- **API Group:** `pvcexplorer.io`
- **Kind:** `PVCExplorerScope`
- **Scope:** Namespaced
- **Status Subresource:** Yes

**Purpose:** Define scoped access policies for PVC exploration

### Commands Reference

```bash
# Local builds
./scripts/build-crds.sh                    # Automated build
kustomize build config/crd                 # Manual CRD build
kustomize build config/default             # Manual full build

# Installation
kubectl apply -f pvc-explorer-crds.yaml    # Install CRDs
kubectl apply -k config/crd                # Install from source

# Verification
kubectl get crd | grep pvcexplorer.io      # Check CRDs
kubectl api-resources | grep pvcexplorer   # Check API resources
kubectl explain pvcexplorer                # View CRD schema

# Troubleshooting
kubectl describe crd pvcexplorers.pvcexplorer.io
kubectl get pvcexplorer -A                 # All PVCExplorers
kubectl get events -n pvc-explorer-system  # Controller events

# Release
git tag v1.0.0                             # Create tag
git push origin v1.0.0                     # Push tag (triggers workflow)
```

### Release Checklist

- [ ] Bump version in code/docs
- [ ] Update CHANGELOG
- [ ] Create git tag: `git tag v1.0.0`
- [ ] Push tag: `git push origin v1.0.0`
- [ ] Wait for GitHub Actions workflow to complete
- [ ] Verify release artifacts exist on GitHub
- [ ] Test installation with: `kubectl apply -f <release-url>`
- [ ] Announce release in documentation/changelog

---

## Related Documentation

- [CRD Installation Guide](./crd-installation.md) - User-facing installation instructions
- [Getting Started Guide](../getting-started.md) - First steps with PVC-Explorer
- [Architecture Guide](../architecture.md) - How PVC-Explorer works
- [Kubernetes CRD Documentation](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)
- [Kustomize Documentation](https://kustomize.io/)
