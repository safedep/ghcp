# GitHub Comments Proxy
GitHub Comments Proxy Service implementing [API](https://buf.build/safedep/api/docs/main:safedep.services.ghcp.v1)

## Background

Proxy Service to allow GitHub Actions safely comment on a PR even when invoked from a forked repository. See 
[GITHUB_TOKEN permissions](https://docs.github.com/en/actions/security-for-github-actions/security-guides/automatic-token-authentication#permissions-for-the-github_token). This service uses a pre-configured `bot`
user account with `GITHUB_TOKEN` to proxy comments from GitHub Actions after appropriate authentication.

### Authentication

Any request to the proxy service must be authenticated to prevent misuse in spamming arbitrary
repositories using this service and its bot user account.

There are multiple methods to authenticate a request that we can consider

1. [GitHub Workload Identity](https://docs.github.com/en/actions/security-for-github-actions/security-hardening-your-deployments/about-security-hardening-with-openid-connect)
2. [GitHub Temporary Token](https://docs.github.com/en/actions/security-for-github-actions/security-guides/automatic-token-authentication)
3. Custom Repository Verification

While [1] seemed like an appropriate solution, unfortunately it is not available when GitHub Actions workflows are executed
from a forked repository due to security reason. It suffers from the same limitation (or security hardening) as `$GITHUB_TOKEN`.

We can leverage the read-only `$GITHUB_TOKEN` to verify the identity of the caller before authorizing the request. However,
we avoid exposing a secret, even though short-lived outside an user's GitHub environment.

We settle for a simple verification of:

- Pre-existing file path in the repository (e.g. `/.github/workflows/vet-ci.yml`)
- Regular expression matching the file content

While this approach appears naive, it is practically useful because we want [vet](https://github.com/safedep/vet)
users, especially open source maintainers to be able to use `vet` especially with PRs from forked repositories.

## Hosted API

A publicly accessible version of the API is hosted at `https://ghcp-integrations.safedep.io`. The API is
authenticated using any of:

1. GitHub Workload Identity token
2. GitHub Repository Verification for `/.github/workflows/vet-ci.yml` as per [vet-action](https://github.com/safedep/vet-action)

### Token Requirements

- GitHub Workload Identity token must be present in `Authorization` header
- GitHub Workload Identity token must have `audience` set to `safedep-ghcp`

## References

- [vet - Policy Driven vetting of OSS Components](https://github.com/safedep/vet)
- [vet-action - GitHub Action for vet](https://github.com/safedep/vet-action)
- [GITHUB_TOKEN permissions](https://docs.github.com/en/actions/security-for-github-actions/security-guides/automatic-token-authentication#permissions-for-the-github_token)
