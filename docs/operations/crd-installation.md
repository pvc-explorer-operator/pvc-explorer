# CRD Installation Guide

This guide explains how to install PVC-Explorer Custom Resource Definitions (CRDs) and deploy the operator.

## What are CRDs?

Custom Resource Definitions (CRDs) extend Kubernetes by defining new resource types. PVC-Explorer provides two CRDs:

- **PVCExplorer**: Manages individual PVC browsing sessions
- **PVCExplorerScope**: Defines scoped access to PVCs

## Installation Methods

### Option 1: Full Installation (Recommended)

Installs CRDs + RBAC + Controller in one command:

```bash
# Replace vX.Y.Z with the desired release version
kubectl apply -f https://github.com/pvc-explorer-operator/pvc-explorer/releases/download/vX.Y.Z/pvc-explorer-install.yaml
