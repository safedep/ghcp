"use strict";

import { createGrpcTransport } from "@connectrpc/connect-node";
import { createClient } from "@connectrpc/connect";
import { GitHubCommentsProxyService } from "@buf/safedep_api.bufbuild_es/safedep/services/ghcp/v1/ghcp_pb";

const apiBaseUrl = "https://ghcp-integrations.safedep.io";

function createGithubCommentsProxyServiceClient() {
  const transport = createGrpcTransport({ baseUrl: apiBaseUrl });
  return createClient(GitHubCommentsProxyService, transport);
}

async function main() {
  const client = createGithubCommentsProxyServiceClient();
  const response = await client.createPullRequestComment({
    owner: "safedep",
    repo: "ghcp",
    prNumber: "1",
    body: "Hello, world!",
  });

  console.log(JSON.stringify(response, null, 2));
}

main();
