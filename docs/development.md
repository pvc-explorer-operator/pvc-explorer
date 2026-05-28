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
