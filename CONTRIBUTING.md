# Contributing to pvc-explorer

Thank you for taking the time to contribute! Every bug report, feature idea, and code improvement makes this project better for everyone.

## Table of contents

- [Code of Conduct](#code-of-conduct)
- [Getting started](#getting-started)
- [How to report a bug](#how-to-report-a-bug)
- [How to suggest a feature](#how-to-suggest-a-feature)
- [How to submit a pull request](#how-to-submit-a-pull-request)
- [Development setup](#development-setup)
- [Commit style](#commit-style)

---

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating you agree to abide by its terms. Please report unacceptable behaviour to the maintainers via a private GitHub message.

---

## Getting started

Not sure where to start? Look for issues labelled **`good first issue`** — these are intentionally scoped to be approachable without deep knowledge of the codebase.

Maintainer roles, responsibilities, and sensitive-access scope are documented in [MAINTAINERS.md](MAINTAINERS.md).

For larger changes, **open an issue first** to discuss the idea before writing code. This avoids wasted effort if the direction doesn't fit the project's scope.

---

## How to report a bug

Use the **Bug report** issue template. Please include:

- What you did
- What you expected to happen
- What actually happened
- Your Kubernetes version, controller version, and how you deployed (kind / in-cluster)

Security vulnerabilities should **not** be reported as public issues — see [SECURITY.md](SECURITY.md).

---

## How to suggest a feature

Use the **Feature request** issue template. Explain the problem you're trying to solve, not just the solution you have in mind. That helps us understand whether it fits the project's scope and discuss alternatives.

---

## How to submit a pull request

1. Fork the repo and create a branch from `main`.
2. Make your changes. If you're fixing a bug, add a test that would have caught it.
3. Run the full test suite locally: `make test`
4. Run the linter: `make lint-fix`
5. If you changed `api/v1alpha1/*_types.go`, run `make manifests generate`.
6. Open a pull request against `main`. Fill in the PR template.

The `main` branch is protected by repository rulesets. Merges require the required CI checks (`Lint`, `Test`, `DCO`, `Trivy SCA`) and at least one approving review (including CODEOWNERS review when applicable). Force pushes and branch deletion are blocked on protected branches. See [docs/operations/branch-protection.md](docs/operations/branch-protection.md).

A maintainer will review within a reasonable time. If you haven't heard back in a week, feel free to ping the thread.

---

## Development setup

**Prerequisites:** Go 1.24+, Node 22+, Docker, [kind](https://kind.sigs.k8s.io/).

```bash
# Clone and install dependencies
git clone https://github.com/pvc-explorer-operator/pvc-explorer.git
cd pvc-explorer

# Run the UI dev server (no cluster needed — uses mock data)
cd ui && npm install && npm run dev

# Run all unit tests
make test

# Run the controller locally against the current kubeconfig context
make run

# Spin up a local kind cluster with everything deployed
./kind/setup.sh
```

See [README.md](README.md) for the full development and release workflow.

---

## Developer Certificate of Origin

By contributing to this project, you certify that:

1. You created the code you're submitting, or you have the right to contribute it.
2. You understand this project and its code will be used under the Apache 2.0 License.
3. You understand your contributions are made under the same license.

We enforce the Developer Certificate of Origin (DCO) on all commits. Sign each commit with:

```bash
git commit -s -m "Your commit message"
```

The `-s` flag adds a `Signed-off-by` line to your commit message. This indicates your agreement with the DCO.

---

## Commit style

We use [Conventional Commits](https://www.conventionalcommits.org/):

```text
feat: add idle timeout override per PVCExplorer
fix: handle RWO PVC on multi-node consumers correctly
docs: clarify mount policy matrix
chore: bump controller-runtime to v0.24
```

This keeps the auto-generated release notes readable.
