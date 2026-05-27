# Security

This page summarizes security guidance for PVC Explorer operator deployments.

## Primary references

- Cluster and platform controls: [Compliance and Security](/compliance-security)
- Operational guardrails: [Operations](/operations)
- Architecture context: [Architecture](/architecture)

## Security checklist

- Use least-privilege RBAC for controller and users.
- Restrict access to the UI ingress and API endpoints.
- Use trusted container image sources and pinned versions.
- Enable audit logging and monitor auth events.
- Review scope boundaries before granting access.
