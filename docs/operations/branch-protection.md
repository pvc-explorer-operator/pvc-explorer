# Branch Protection and Rulesets

This project uses GitHub repository rulesets (not classic branch protection) to protect `main` and release branches.

## Active Ruleset

- **Ruleset name**: Main/release branch protection (ruleset)
- **Target branches**:
  - `refs/heads/main`
  - `refs/heads/release/*`
- **Enforcement**: Active

## Required Checks

The following status checks are required before merge:

- `Lint`
- `Test`
- `DCO`
- `Trivy SCA`
- `OSPS Assessment`

The ruleset requires branches to be up to date with the base branch before merge.

## Pull Request Requirements

Direct pushes to protected branches are not the normal workflow. Changes are expected through pull requests with:

- At least 1 approving review
- Code owner review required
- Last push approval required
- Stale approvals dismissed on new commits
- Review threads resolved before merge

## Additional Protections

- Force pushes are blocked
- Branch deletion is blocked
- No bypass actors are configured

## Repository Administration

Maintainers with repository admin access manage this ruleset. See [MAINTAINERS.md](../../MAINTAINERS.md) for role and access scope.

## Audit and Verification

Use the GitHub CLI to inspect the active ruleset:

```bash
gh api repos/pvc-explorer-operator/pvc-explorer/rulesets/17040093
```

To list all rulesets:

```bash
gh api repos/pvc-explorer-operator/pvc-explorer/rulesets
```
