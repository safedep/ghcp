# GitHub Comments Proxy
GitHub Comments Proxy Service implementing [API](https://buf.build/safedep/api/docs/main:safedep.services.ghcp.v1)

## Background

Proxy Service to allow GitHub Actions safely comment on a PR even when invoked from a forked repository. See [GITHUB_TOKEN permissions](https://docs.github.com/en/actions/security-for-github-actions/security-guides/automatic-token-authentication#permissions-for-the-github_token).

## Hosted API

A publicly accessible version of the API is hosted at `https://ghcp.integrations.safedep.io`.

### Token Requirements

- GitHub Workload Identity token must be present in `Authorization` header
- GitHub Workload Identity token must have `audience` set to `safedep-ghcp`
