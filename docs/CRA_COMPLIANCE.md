# Cyber Resilience Act (CRA) Compliance

This document explains how pvc-explorer addresses the European Union's Cyber Resilience Act (CRA) requirements for open-source software.

## Overview

The **Cyber Resilience Act (CRA)**, effective from 2025, establishes cybersecurity requirements for digital products in the EU. It includes special provisions for free and open-source software (FOSS) to recognize their unique development models and societal importance.

**Relevant resources:**
- [EU CRA Summary](https://digital-strategy.ec.europa.eu/en/policies/cra-summary)
- [CRA Legal Text](https://eur-lex.europa.eu/legal-content/EN/TXT/?uri=CELEX%3A32024R2847)
- [OpenSSF CRA Course](https://openssf.org/blog/2025/06/16/cra-ready-how-open-source-projects-can-prepare-for-the-eu-cyber-resilience-act/)

## Open-Source Software Steward Status

**pvc-explorer Status**: Development

If pvc-explorer transitions to a commercial FOSS offering (offered on the market with sustained support), it may qualify as an **open-source software steward** under CRA Article 24.

### Definition of a Steward

A legal person that:

1. Does not place the software on the market themselves (not monetized by the project)
2. Provides systematic support on a sustained basis for FOSS development
3. Plays a main role in ensuring the viability of products intended for commercial activities

**Current Status**: pvc-explorer is an unpaid open-source project. No steward obligations apply yet.

## When Steward Obligations Would Apply

If pvc-explorer becomes:

- ✅ Offered as a service with SLA/support
- ✅ Bundled with commercial products
- ✅ Supported by a legal entity as a sustained service

Then the following obligations would apply:

### Article 24 Obligations for Stewards

#### 1. Cybersecurity Policy

A documented policy fostering secure development practices:
- Secure coding guidelines ([docs/development.md](development.md))
- Vulnerability management process ([SECURITY.md](../SECURITY.md))
- Security review practices for contributions
- Dependency management and license compliance ([docs/LICENSE_COMPLIANCE.md](LICENSE_COMPLIANCE.md))

**pvc-explorer currently has:**
- ✅ Contribution guidelines ([CONTRIBUTING.md](../CONTRIBUTING.md)) with security focus
- ✅ Vulnerability disclosure policy ([SECURITY.md](../SECURITY.md))
- ✅ License compliance tracking ([docs/LICENSE_COMPLIANCE.md](LICENSE_COMPLIANCE.md))
- ✅ Code review requirements (PR-based development)

#### 2. Vulnerability Handling

- Establish a process for reporting and handling vulnerabilities
- Provide security updates timely
- Coordinate disclosure of actively exploited vulnerabilities

**pvc-explorer currently has:**
 - ✅ Private vulnerability reporting ([SECURITY.md](../SECURITY.md))
 - ⚠️ Best-effort response to vulnerability reports (no guaranteed timeline; this is a no-cost, community-supported project with no warranty or responsibility)
 - ✅ Security advisory process
 - ✅ Explicit coordination of releases for fixes (as maintainers are able)

#### 3. Market Surveillance Cooperation

- Cooperate with market surveillance authorities
- Report on the cybersecurity status of the software
- No administrative fines apply to stewards (Article 64(10))

**pvc-explorer currently has:**
- ✅ Maintainer contact information (GitHub issues/discussions)
- ✅ Public issue tracking for known issues
- ✅ Release notes documenting security fixes

#### 4. Active Vulnerability Reporting

- Report actively exploited vulnerabilities to CISA (if steward)
- Coordinate with other entities to mitigate impact

**pvc-explorer currently has:**
- ✅ Process for receiving vulnerability reports
- ✅ Ability to coordinate with ecosystem partners
- (Reporting to CISA would be initiated if/when vulnerability becomes known)

## Current Compliance Posture

| Requirement                 | Status   | Evidence                                            |
| --------------------------- | -------- | --------------------------------------------------- |
| Contributor agreements      | ✅ Ready | [CONTRIBUTING.md](../CONTRIBUTING.md) with DCO      |
| Cybersecurity policy        | ✅ Ready | [SECURITY.md](../SECURITY.md)                       |
| License compliance tracking | ✅ Ready | [docs/LICENSE_COMPLIANCE.md](LICENSE_COMPLIANCE.md) |
| SBOM generation             | ✅ Ready | `make sbom` target                                  |
| Transparent development     | ✅ Ready | Public GitHub repository                            |
| Vulnerability handling      | ✅ Ready | [SECURITY.md](../SECURITY.md)                       |

## Preparing for Steward Status

If pvc-explorer transitions to a steward model in the future:

1. **Formalize support model**: Document response times, SLA, support channels
2. **Enhance testing**: Increase test coverage, add fuzzing for security-critical paths
3. **Dependency management**: Actively monitor and update dependencies for vulnerabilities
4. **Security scanning**: Integrate SAST/DAST tools in CI/CD pipeline
5. **Incident response**: Define and document IR process
6. **Audit trail**: Maintain records of security decisions and vulnerability handling

## Automated Vulnerability Scanning (Grype)

pvc-explorer uses [Grype](https://github.com/anchore/grype) to scan for vulnerabilities in dependencies and container images as part of its CRA compliance process.

- Grype scans are run on the built container image and/or SBOM to ensure all shipped dependencies are checked for known vulnerabilities.
- Vulnerability reports are generated and stored as release artifacts and/or in CI logs.
- All vulnerabilities with available fixes are tracked and remediated before release.
- The process is automated in CI/CD and documented in the [MIGRATION_CHECKLIST.md](MIGRATION_CHECKLIST.md).

**Known limitations:**
- Some vulnerabilities may be reported due to upstream or ecosystem issues (e.g., incomplete SBOM data, false positives). These are documented and tracked for due diligence.
- The project monitors CNCF/OpenSSF and upstream projects for improvements in SBOM and vulnerability reporting.

**Example Grype usage:**
```bash
# Scan built image
grype ghcr.io/<org>/pvc-explorer:latest
# Scan SBOM
grype sbom:dist/sbom.cyclonedx.json
```

See also: [docs/MIGRATION_CHECKLIST.md](MIGRATION_CHECKLIST.md) for the full compliance workflow.

## Resources for FOSS Projects

**Helpful community initiatives:**

- [Open Regulatory Compliance WG (ORCWG) on CRA](https://orcwg.org/cra/)
- [OpenSSF Blog on CRA](https://openssf.org/blog/2025/06/16/cra-ready-how-open-source-projects-can-prepare-for-the-eu-cyber-resilience-act/)
- [CISA Vulnerability Coordination](https://www.cisa.gov/coordinated-vulnerability-disclosure-practice)

## Contact

For questions regarding CRA compliance or security practices:

1. **Security vulnerabilities**: [SECURITY.md](../SECURITY.md) (private reporting)
2. **General questions**: Open a GitHub discussion or issue
3. **Maintainers**: Contact via GitHub or repository contact info

## Appendix: CRA Terminology

- **Product with Digital Elements**: Any product that contains digital components (hardware or software)
- **Free and Open Source Software (FOSS)**: Software distributed under an OSI-approved license
- **Open-Source Software Steward**: Legal person providing sustained support for FOSS
- **Actively Exploited Vulnerability**: Vulnerability with known public exploit or active attacks
- **Market Surveillance**: EU authority oversight of product safety and compliance

---

**Last Updated**: May 2026

For the latest CRA guidance, see [digital-strategy.ec.europa.eu/en/policies/cra](https://digital-strategy.ec.europa.eu/en/policies/cra)
