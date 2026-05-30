# Development

## Commands

```bash
make test
make manifests
make generate
make lint-fix
```

Always run `make manifests generate` after editing `api/v1alpha1/*_types.go`.

## UI development

The Vite dev server includes a mock plugin so you can work on the UI without a running cluster.

```bash
cd ui
npm install
npm run dev
```

The mock plugin only runs during `npm run dev`.

## UI docs map

- UI local run and auth bypass: [`ui/README.md`](../ui/README.md)
- UI documentation landing page: [`docs/ui/index.md`](ui/index.md)
- UI accessibility guide: [`docs/ui/accessibility.md`](ui/accessibility.md)

## Versioning

The binary version is injected via Go linker flags from the build system and exposed at `GET /api/version`.

Local builds use `git describe --tags --always --dirty`, which makes the version reflect the exact checkout you tested.

## CI metadata validation guardrails

GitHub Actions event metadata is treated as untrusted input in shell steps.

- Validate branch/tag/repository-derived values against an allowlist regex before use.
- Use validated outputs (`$GITHUB_OUTPUT`) in later commands instead of raw event fields.
- Quote shell variables when passing them to commands, paths, or API routes.

Current examples in this repository:

- `.github/workflows/oci-image.yml` validates release tags and package path components before release/package API operations.
- `.github/workflows/release-crds.yaml` validates tag metadata before constructing release assets.
- `.github/workflows/docs-pages.yml` validates repository name before composing `DOCS_BASE_PATH`.

## CI credential boundaries (OSPS-BR-01.03)

Fork and pull-request workflows are treated as untrusted execution contexts.

- PR workflows use read-only permissions where possible (`contents: read`).
- PR workflows do not use privileged repository secrets.
- PR workflows use the built-in `${{ github.token }}` only for read operations required by tooling.

Privileged credentials and write-capability jobs are restricted to trusted contexts:

- Release and publishing workflows run on protected events (for example version tag pushes, scheduled jobs, and maintainer-triggered `workflow_dispatch`).
- Jobs that write packages/releases or use OIDC/signing are isolated to those trusted workflows.
- Branch protections require reviews and required checks before merge.
