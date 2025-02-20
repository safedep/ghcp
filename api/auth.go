package api

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"connectrpc.com/connect"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt"
	"github.com/safedep/dry/log"
	"github.com/safedep/ghcp/pkg/adapters/github"
	"github.com/safedep/ghcp/pkg/gh"
)

type AuthenticationInterceptorConfig struct {
	MockAuthentication bool
}

type authenticationInterceptor struct {
	config   AuthenticationInterceptorConfig
	provider *oidc.Provider
}

// AuthInterceptor is a Connect interceptor that authenticates requests
// using GitHub Workload Identity Token.
func NewAuthenticationInterceptor(config AuthenticationInterceptorConfig) (connect.Interceptor, error) {
	provider, err := oidc.NewProvider(context.Background(), "https://token.actions.githubusercontent.com")
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider for GitHub Workload Identity: %w", err)
	}

	return &authenticationInterceptor{config: config, provider: provider}, nil
}

func (i *authenticationInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		if i.config.MockAuthentication {
			return next(ctx, req)
		}

		authHeader := req.Header().Get("authorization")
		if authHeader == "" {
			return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("token is required"))
		}

		authHeader = strings.TrimPrefix(authHeader, "Bearer ")
		if authHeader == "" {
			return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("token is missing"))
		}

		var tokenContext gh.GitHubTokenContext
		var err error

		if i.isPAT(authHeader) {
			tokenContext, err = i.authenticateUsingPAT(ctx, authHeader)
		} else {
			tokenContext, err = i.authenticateUsingJWT(ctx, authHeader)
		}

		if err != nil {
			return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("failed to authenticate: %w", err))
		}

		ctx = gh.InjectGitHubTokenContext(ctx, tokenContext)
		return next(ctx, req)
	}
}

func (i *authenticationInterceptor) isPAT(token string) bool {
	// https://github.blog/changelog/2021-03-31-authentication-token-format-updates-are-generally-available/
	var validPatPrefixes = []string{
		"ghp_",
		"gho_",
		"ghu_",
		"ghs_",
	}

	for _, prefix := range validPatPrefixes {
		if strings.HasPrefix(token, prefix) {
			return true
		}
	}

	return false
}

// authenticateUsingPAT authenticates the PAT token and returns the internal GitHub token context
func (i *authenticationInterceptor) authenticateUsingPAT(ctx context.Context, token string) (gh.GitHubTokenContext, error) {
	log.Debugf("Authenticating using GITHUB_TOKEN")

	adapter, err := github.NewGitHubAdapter(github.GitHubAdapterConfig{
		Token: token,
	})
	if err != nil {
		return gh.GitHubTokenContext{}, fmt.Errorf("failed to create GitHub adapter: %w", err)
	}

	userInfo, err := adapter.GetTokenUser(ctx, token)
	if err != nil {
		return gh.GitHubTokenContext{}, fmt.Errorf("failed to get token user: %w", err)
	}

	log.Debugf("Token user: %+v", userInfo)

	var tokenContext gh.GitHubTokenContext

	if userInfo.GetLogin() != "" {
		tokenContext.Actor = userInfo.GetLogin()
	}

	log.Debugf("Token context: %+v", tokenContext)

	return tokenContext, nil
}

// authenticateUsingJWT authenticates the OIDC token and returns the internal GitHub token context
func (i *authenticationInterceptor) authenticateUsingJWT(ctx context.Context, authHeader string) (gh.GitHubTokenContext, error) {
	log.Debugf("Authenticating using JWT")

	var tokenContext gh.GitHubTokenContext

	// Authenticate the OIDC token
	verifier := i.provider.Verifier(&oidc.Config{SkipClientIDCheck: true})
	_, err := verifier.Verify(ctx, authHeader)
	if err != nil {
		return tokenContext, connect.NewError(connect.CodeUnauthenticated, errors.New("token verification failed"))
	}

	// We need to re-parse the token to get the GitHub specific claims
	// We don't need to validate the token, we just need to parse it
	parser := &jwt.Parser{}
	claims := jwt.MapClaims{}
	_, _, err = parser.ParseUnverified(authHeader, claims)
	if err != nil {
		return tokenContext, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("token parsing failed: %w", err))
	}

	log.Debugf("Token claims: %+v", claims)

	if s, ok := claims["iss"].(string); ok {
		tokenContext.Issuer = s
	}

	if s, ok := claims["sub"].(string); ok {
		tokenContext.Subject = s
	}

	if s, ok := claims["environment"].(string); ok {
		tokenContext.Environment = s
	}

	if s, ok := claims["repository"].(string); ok {
		tokenContext.Repository = s
	}

	if s, ok := claims["repository_owner"].(string); ok {
		tokenContext.RepositoryOwner = s
	}

	if s, ok := claims["repository_id"].(string); ok {
		tokenContext.RepositoryID = s
	}

	if s, ok := claims["repository_owner_id"].(string); ok {
		tokenContext.RepositoryOwnerID = s
	}

	if s, ok := claims["repository_visibility"].(string); ok {
		tokenContext.RepositoryVisibility = s
	}

	if s, ok := claims["ref"].(string); ok {
		tokenContext.Ref = s
	}

	if s, ok := claims["run_id"].(string); ok {
		tokenContext.RunID = s
	}

	if s, ok := claims["run_number"].(string); ok {
		tokenContext.RunNumber = s
	}

	if s, ok := claims["run_attempt"].(string); ok {
		tokenContext.RunAttempt = s
	}

	if s, ok := claims["runner_environment"].(string); ok {
		tokenContext.RunnerEnvironment = s
	}

	if s, ok := claims["actor"].(string); ok {
		tokenContext.Actor = s
	}

	if s, ok := claims["workflow"].(string); ok {
		tokenContext.Workflow = s
	}

	if s, ok := claims["workflow_ref"].(string); ok {
		tokenContext.WorkflowRef = s
	}

	log.Debugf("Token context: %+v", tokenContext)

	return tokenContext, nil
}

func (i *authenticationInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return func(ctx context.Context, spec connect.Spec) connect.StreamingClientConn {
		return nil
	}
}

func (i *authenticationInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, stream connect.StreamingHandlerConn) error {
		return fmt.Errorf("not implemented")
	}
}
