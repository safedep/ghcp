"use strict";

import { createGrpcTransport } from "@connectrpc/connect-node";
import { createClient } from "@connectrpc/connect";
import { GitHubCommentsProxyService } from "@buf/safedep_api.bufbuild_es/safedep/services/ghcp/v1/ghcp_pb";

const apiBaseUrl = "https://ghcp-integrations.safedep.io";

function authenticationInterceptor(token: string) {
  return (next: any) => async (req: any) => {
    req.header.set("authorization", `Bearer ${token}`);
    return await next(req);
  };
}

function createGithubCommentsProxyServiceClient(token: string) {
  const transport = createGrpcTransport({
    baseUrl: apiBaseUrl,
    interceptors: [authenticationInterceptor(token)],
  });

  return createClient(GitHubCommentsProxyService, transport);
}

async function main() {
  const githubToken = process.env.GITHUB_TOKEN;
  if (!githubToken) {
    throw new Error("GITHUB_TOKEN is not set");
  }

  const githubPullRequestNumber = process.env.GITHUB_PULL_REQUEST_NUMBER;
  if (!githubPullRequestNumber) {
    throw new Error("GITHUB_PULL_REQUEST_NUMBER is not set");
  }

  const client = createGithubCommentsProxyServiceClient(githubToken);
  const response = await client.createPullRequestComment({
    owner: "safedep",
    repo: "ghcp",
    prNumber: githubPullRequestNumber,
    body: "Hello, world!",
  });

  console.log(JSON.stringify(response, null, 2));
}

main();
