package api

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"connectrpc.com/connect"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt"
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
			return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("unauthenticated"))
		}

		authHeader = strings.TrimPrefix(authHeader, "Bearer ")
		if authHeader == "" {
			return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("unauthenticated"))
		}

		// Authenticate the OIDC token
		verifier := i.provider.Verifier(&oidc.Config{SkipClientIDCheck: true})
		_, err := verifier.Verify(ctx, authHeader)
		if err != nil {
			return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("token verification failed"))
		}

		// We need to re-parse the token to get the GitHub specific claims
		// We don't need to validate the token, we just need to parse it
		parser := &jwt.Parser{}
		claims := jwt.MapClaims{}
		_, _, err = parser.ParseUnverified(authHeader, claims)
		if err != nil {
			return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("token parsing failed"))
		}

		// Extract and build token context
		var tokenContext gh.GitHubTokenContext

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

		// Inject token context into the context
		ctx = gh.InjectGitHubTokenContext(ctx, tokenContext)

		// Call next handler
		return next(ctx, req)
	}
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
