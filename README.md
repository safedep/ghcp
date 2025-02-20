# GitHub Comments Proxy
GitHub Comments Proxy Service implementing [API](https://buf.build/safedep/api/docs/main:safedep.services.ghcp.v1). This service
is built to help GitHub Actions developers to comment on a PR even when invoked from a forked repository.

## TL;DR

```bash
curl -X POST \
  https://ghcp-integrations.safedep.io/safedep.services.ghcp.v1.GitHubCommentsProxyService/CreatePullRequestComment \
  -H "Authorization: Bearer $GITHUB_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"owner": "safedep", "repo": "ghcp", "pr_number": "1", "body": "Hello, world!"}'
```

For the request to be successful, the following conditions must be met:

- `$GITHUB_TOKEN` is a temporary GitHub Action Token (NOT user PAT)
- `$GITHUB_TOKEN` has access to the requested repository
- The requested PR is in `open` state

## Hosted API

A publicly accessible version of the API is hosted at `https://ghcp-integrations.safedep.io`. The API is
authenticated using any of:

1. GitHub [Workload Identity token](https://docs.github.com/en/actions/security-for-github-actions/security-hardening-your-deployments/about-security-hardening-with-openid-connect)
2. GitHub [Actions Token](https://docs.github.com/en/actions/security-for-github-actions/security-guides/automatic-token-authentication#permissions-for-the-github_token) (`$GITHUB_TOKEN`)

### Limits

- Maximum 3 comments per PR
- Unlimited comment updates using `tag` subject to GitHub API rate limits

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
we want to avoid exposing a secret, even though short-lived outside an user's GitHub environment.

Verification of the repository is done by checking for the existence of a pre-existing file path in the repository.

- Pre-existing file path in the repository (e.g. `/.github/workflows/vet-ci.yml`)
- Regular expression matching the file content

However, this approach is vulnerable to spamming all existing users of [vet-action](https://github.com/safedep/vet-action).
We settled for a restricted used of `$GITHUB_TOKEN` with following verification:

- Verify `$GITHUB_TOKEN` is present in the request and it is a temporary GitHub Action Token (NOT user PAT)
- Verify `$GITHUB_TOKEN` has access to the requested repository
- Verify the requested PR is in `open` state

Only then the service accepts the request and proxies the comment to the target repository.

## Security

Send all security reports to `security@safedep.io`.

## References

- [vet - Policy Driven vetting of OSS Components](https://github.com/safedep/vet)
- [vet-action - GitHub Action for vet](https://github.com/safedep/vet-action)
- [GITHUB_TOKEN permissions](https://docs.github.com/en/actions/security-for-github-actions/security-guides/automatic-token-authentication#permissions-for-the-github_token)
