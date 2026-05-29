# Maintainers

This file documents project members with access to sensitive resources and their responsibilities.

## Roles

### Maintainer

Maintainers are trusted project members who can approve and merge changes, manage releases, and administer repository security settings.

Responsibilities:

- Review and merge pull requests.
- Manage releases, tags, and published artifacts.
- Maintain CI/CD, security checks, and repository settings.
- Handle vulnerability reports according to `SECURITY.md`.

### Contributor

Contributors can propose changes through pull requests but do not have maintainer-level administrative access.

Responsibilities:

- Submit issues and pull requests.
- Follow contribution and security policies.
- Participate in code review discussions.

## Current Maintainers

| Name | GitHub | Role | Sensitive resource access |
| --- | --- | --- | --- |
| Ricardo Leal | [@ricardoleal](https://github.com/ricardoleal) | Maintainer | Repository admin settings, branch protection and rulesets, GitHub Actions settings and secrets, release/tag publishing |

## Sensitive Resources

For this project, sensitive resources include:

- GitHub repository admin settings (branch protection, rulesets, collaborator/admin access).
- GitHub Actions secrets and workflow permissions.
- Release pipelines, release tags, and published OCI artifacts.
- Private vulnerability reports in GitHub Security.

## Process to Add a New Maintainer

1. Open a pull request updating this file with the new maintainer's name, GitHub handle, and role.
2. Include the reason for access and which sensitive resources are needed.
3. Obtain approval from at least one existing maintainer.
4. After merge, apply the corresponding GitHub permission changes.
5. Confirm the maintainer has read `CONTRIBUTING.md` and `SECURITY.md`.

## Process to Remove or Change Maintainer Access

1. Open a pull request updating this file to remove or modify maintainer access.
2. Apply matching permission changes in GitHub settings immediately after merge.
3. Rotate or revoke credentials/tokens if the role change requires it.
