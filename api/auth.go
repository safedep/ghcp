package api

import (
	"context"
	"errors"
	"fmt"

	"connectrpc.com/connect"
)

type AuthenticationInterceptorConfig struct {
	MockAuthentication bool
}

type authenticationInterceptor struct {
	config AuthenticationInterceptorConfig
}

// AuthInterceptor is a Connect interceptor that authenticates requests
// using GitHub Workload Identity Token.
func NewAuthenticationInterceptor(config AuthenticationInterceptorConfig) (connect.Interceptor, error) {
	return &authenticationInterceptor{config: config}, nil
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

		// Authenticate the OIDC token
		// Extract and build token context
		// Inject token context into the context
		// Call next handler

		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("unauthenticated"))
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
