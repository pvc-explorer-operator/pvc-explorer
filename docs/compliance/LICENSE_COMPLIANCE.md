# License Compliance

This document describes how pvc-explorer manages third-party dependencies and ensures license compliance.

## License Overview

**pvc-explorer** is licensed under the **Apache License 2.0**, which:

- Grants unlimited permissions for use, modification, and distribution
- Includes an explicit patent grant from contributors
- Requires preservation of copyright notices and license text
- Provides no warranty or liability protection to users

See [LICENSE](../../LICENSE) for the full text.

## Dependency License Requirements

All dependencies must be compatible with Apache 2.0. Acceptable license categories:

| Category             | Examples                           | Compatible? |
| -------------------- | ---------------------------------- | ----------- |
| **Permissive**       | MIT, BSD-3-Clause, Apache 2.0, ISC | ✅ Yes      |
| **Source-Available** | BSL, SSPL                          | ❌ No       |
| **Strong Copyleft**  | GPLv2, GPLv3, AGPL                 | ❌ No       |
| **Weak Copyleft**    | MPL 2.0, LGPL 2.1+, CDDL           | ✅ Yes      |

## Go Dependencies

### Checking Licenses

List all direct dependencies:
```bash
go mod graph
```

View license text for a specific package:
```bash
go-licenses csv ./...      # CSV report of all licenses
go-licenses report ./...   # Detailed license report
go-licenses read <PACKAGE> # View license text
```

### Compliance Strategy

- All dependencies are listed in `go.mod`
- Indirect dependencies are in `go.sum`
- CI/CD validates that no incompatible licenses are introduced
- License changes in dependencies are reviewed during dependency updates

### Tools Used

- **go-licenses**: Audits and reports on Go package licenses
  ```bash
  go install github.com/google/go-licenses@latest
  ```

## npm Dependencies (UI)

### Checking Licenses

List all dependencies and their licenses:
```bash
cd ui && npm ls --depth=0
```

View a specific package's license:
```bash
cd ui && npm view <PACKAGE> license
```

### Compliance Strategy

- All npm packages are listed in `ui/package.json`
- Lock file: `ui/package-lock.json` pins exact versions
- Development runs `npm audit` to catch security and license issues
- CI/CD validates compatibility before merging

## Software Bill of Materials (SBOM)

An SBOM is generated for each release and includes:

- **CycloneDX format** (`sbom.cyclonedx.json`) — standardized for enterprise compliance
- **SPDX format** (`sbom.spdx.json`) — standardized for license compliance

### Generate SBOM Locally

```bash
make sbom
```

Outputs:
- `dist/sbom.cyclonedx.json` — CycloneDX JSON (NIST/enterprise standard)
- `dist/sbom.spdx.json` — SPDX JSON (raw Syft output)
- `dist/sbom.normalized.spdx.json` — SPDX JSON with normalized license conclusions for audit tooling

### NOASSERTION Interpretation

Some scanners populate SPDX fields such as `licenseConcluded`, `supplier`, or `copyrightText`
with `NOASSERTION` when upstream metadata is incomplete. This does not automatically mean the
package uses a non-approved license.

For compliance decisions in this project:

- The raw SPDX output is retained as evidence (`dist/sbom.spdx.json`).
- A normalized SPDX variant is generated for audit compatibility (`dist/sbom.normalized.spdx.json`).
- The authoritative compatibility gate is `go-licenses check` with the allowlist below.

### SBOM in Releases

Each release includes:
- SBOM files in the release artifacts (raw and normalized SPDX)
- Reference in release notes: _"SBOM available in release assets"_
- Can be used by organizations for compliance scanning

## Third-Party Material

### Non-Source Code Content

The UI includes third-party assets:
- **Icons**: PrimeIcons (MIT) — distributed via `@primevue/themes`
- **CSS Framework**: PrimeVue (MIT) — component library

All are MIT licensed (compatible with Apache 2.0).

### Images and Logos

- Logo and icons: Internal assets
- Documentation images: Internal assets

## Verified Dependencies (as of May 2026)

### Summary

✅ **100% Compliant** — All 100+ external dependencies are Apache 2.0 compatible.

**License Breakdown:**
- **Apache-2.0**: 40+ (Kubernetes, controller-runtime, OpenTelemetry, Google libraries)
- **MIT**: 20+ (Testing libraries, utilities, logging)
- **BSD-3-Clause**: 15+ (Go standard library forks, Kubernetes)
- **ISC**: 2 (websocket, go-spew)

**Key Dependencies:**
- ✅ `k8s.io/*` → Apache-2.0 (Kubernetes ecosystem)
- ✅ `sigs.k8s.io/controller-runtime` → Apache-2.0 (Kubebuilder core)
- ✅ `github.com/onsi/ginkgo/v2` → MIT (Testing framework)
- ✅ `github.com/prometheus/client_golang` → Apache-2.0 (Metrics)
- ✅ `go.opentelemetry.io/*` → Apache-2.0 (Observability)

**No problematic licenses found:** ❌ No GPLv2, GPLv3, AGPL, BSL, or SSPL dependencies.

---

## Legal Review Checklist

Before each release, verify:

- [ ] All direct dependencies listed in `go.mod` are compatible
- [ ] Run `go-licenses check ./...` with the allowed whitelist (see below)
- [ ] SBOM generated and validated
- [ ] No new strong copyleft dependencies introduced
- [ ] Copyright notices in source files are current
- [ ] LICENSE and NOTICE files are present and accurate
- [ ] Third-party materials are properly attributed

### Dependency Whitelist

Use this command to validate compliance:

```bash
go-licenses check ./... \
  --allowed_licenses=MIT,Apache-2.0,BSD-3-Clause,BSD-2-Clause,ISC,MPL-2.0,LGPL-2.1-or-later,LGPL-3.0-or-later,CDDL-1.0 \
  --ignore github.com/pvc-explorer-operator/pvc-explorer
```

**Explanation:**
- `--allowed_licenses` — Whitelist of compatible licenses
- `--ignore github.com/pvc-explorer-operator/pvc-explorer` — Skip local code (covered by root LICENSE)

## Reporting License Violations

If you believe a dependency violates Apache 2.0 compatibility:

1. Open a **private security advisory** on GitHub
2. Provide:
   - Affected dependency name and version
   - Incompatible license
   - Potential impact

See [SECURITY.md](../../SECURITY.md) for details.

## References

- [Apache License 2.0](https://opensource.org/licenses/Apache-2.0)
- [SPDX License List](https://spdx.org/licenses/)
- [Open Source Guide: Legal](https://opensource.guide/legal/)
- [Google go-licenses](https://github.com/google/go-licenses)
- [CycloneDX Standard](https://cyclonedx.org/)
- [SPDX Standard](https://spdx.dev/)
