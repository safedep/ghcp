package gh

import (
	"context"
	"errors"
)

type githubTokenContextKey struct{}

type TokenType string

const (
	TokenTypeAction           TokenType = "action"
	TokenTypeUser             TokenType = "user"
	TokenTypeWorkloadIdentity TokenType = "workload_identity"
)

// GitHubTokenContext holds information extracted from the GitHub Workload Identity Token
// https://docs.github.com/en/actions/security-for-github-actions/security-hardening-your-deployments/about-security-hardening-with-openid-connect
type GitHubTokenContext struct {
	Subject              string `json:"sub"`
	Issuer               string `json:"iss"`
	Environment          string `json:"environment"`
	Audience             string `json:"aud"`
	Repository           string `json:"repository"`
	RepositoryOwner      string `json:"repository_owner"`
	RepositoryVisibility string `json:"repository_visibility"`
	RepositoryID         string `json:"repository_id"`
	RepositoryOwnerID    string `json:"repository_owner_id"`
	Ref                  string `json:"ref"`
	RunID                string `json:"run_id"`
	RunNumber            string `json:"run_number"`
	RunAttempt           string `json:"run_attempt"`
	RunnerEnvironment    string `json:"runner_environment"`
	Actor                string `json:"actor"`
	Workflow             string `json:"workflow"`
	WorkflowRef          string `json:"workflow_ref"`
	WorkflowSHA          string `json:"workflow_sha"`
	HeadRef              string `json:"head_ref"`
	BaseRef              string `json:"base_ref"`
	RefType              string `json:"ref_type"`
	EventName            string `json:"event_name"`
	JobWorkflowRef       string `json:"job_workflow_ref"`

	// TokenType is the type of token
	TokenType TokenType
}

// Inject GitHub token context into the context
func InjectGitHubTokenContext(ctx context.Context, tokenContext GitHubTokenContext) context.Context {
	return context.WithValue(ctx, githubTokenContextKey{}, tokenContext)
}

// Extract GitHub token context from the context
func ExtractGitHubTokenContext(ctx context.Context) (GitHubTokenContext, error) {
	tokenContext, ok := ctx.Value(githubTokenContextKey{}).(GitHubTokenContext)
	if !ok {
		return GitHubTokenContext{}, errors.New("no GitHub token context found")
	}

	return tokenContext, nil
}

func (t *GitHubTokenContext) SetTokenType(tokenType TokenType) {
	t.TokenType = tokenType
}

func (t GitHubTokenContext) IsActionToken() bool {
	return t.TokenType == TokenTypeAction
}

func (t GitHubTokenContext) IsUserToken() bool {
	return t.TokenType == TokenTypeUser
}

func (t GitHubTokenContext) IsWorkloadIdentityToken() bool {
	return t.TokenType == TokenTypeWorkloadIdentity
}
