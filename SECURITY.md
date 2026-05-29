# Security Policy

## Supported Versions

Security fixes are applied to the latest release on `main`.

| Version        | Supported          |
| -------------- | ------------------ |
| Latest `main`  | :white_check_mark: |
| Older releases | :x:                |

## Reporting a Vulnerability

Please do **not** report security vulnerabilities in public GitHub issues.

Maintainer roles and sensitive-resource access are documented in [MAINTAINERS.md](MAINTAINERS.md).

Use GitHub's private vulnerability reporting:

1. Go to the repository's **Security** tab
2. Click **Report a vulnerability**
3. Provide reproduction steps and impact details

If private reporting is unavailable, contact maintainers via GitHub private message.

We will acknowledge receipt within 72 hours and provide a status update as triage progresses.

## Disclosure Process

- We confirm the issue and assess severity
- We prepare and test a fix
- We coordinate a release
- We publish an advisory with mitigation details

## Preventing Secret Leakage

This repository must not store real secrets or credentials in version control.

Policy:

- Do not commit tokens, passwords, API keys, kubeconfigs, private keys, or session secrets.
- Use placeholders or mock values for examples and local development fixtures.
- Keep operational secrets in environment variables, GitHub secrets, or Kubernetes Secrets, not in tracked files.

Automated enforcement:

- CI runs automated secret scanning on pull requests and pushes to `main`.
- Pull requests that introduce suspected secrets must be fixed before merge.
- If a scan reports a false positive, update the scanner allowlist with a narrow, file-specific exception and explain why it is safe.
