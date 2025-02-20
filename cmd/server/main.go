package server

import (
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	dryhttp "github.com/safedep/dry/adapters/http"
	"github.com/safedep/dry/log"
	"github.com/safedep/dry/obs"
	"github.com/safedep/ghcp/api"
	"github.com/safedep/ghcp/pkg/adapters/github"
	"github.com/safedep/ghcp/services/ghcp"
	"github.com/spf13/cobra"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var (
	serverAddress            string
	serverMockAuthentication bool
	serverMockAuthorization  bool
)

func NewServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "server",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := startServer()
			if err != nil {
				log.Fatalf("failed to start server: %v", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&serverAddress, "address", "127.0.0.1:8000", "address to listen on")
	cmd.Flags().BoolVar(&serverMockAuthentication, "mock-authentication", false, "enable mock authentication")
	cmd.Flags().BoolVar(&serverMockAuthorization, "mock-authorization", false, "enable mock authorization")
	return cmd
}

func startServer() error {
	interceptors, err := buildConnectInterceptors()
	if err != nil {
		return fmt.Errorf("failed to build connect interceptors: %w", err)
	}

	router, err := dryhttp.NewEchoRouter(dryhttp.EchoRouterConfig{
		ServiceName: obs.AppServiceName("ghcp"),
	})
	if err != nil {
		return fmt.Errorf("failed to create echo router: %w", err)
	}

	githubAdapter, err := github.NewGitHubAdapter(github.DefaultGitHubAdapterConfig())
	if err != nil {
		return fmt.Errorf("failed to create github issue adapter: %w", err)
	}

	ghcpServiceConfig := ghcp.DefaultGitHubCommentProxyServiceConfig()
	ghcpServiceConfig.InsecureSkipAuthorization = serverMockAuthorization

	ghcpService, err := ghcp.NewGitHubCommentProxyService(ghcpServiceConfig,
		githubAdapter, githubAdapter)
	if err != nil {
		return fmt.Errorf("failed to create ghcp service: %w", err)
	}

	apiHandler, err := api.NewGhcpServiceHandler(ghcpService)
	if err != nil {
		return fmt.Errorf("failed to create ghcp service handler: %w", err)
	}

	err = registerService(router, apiHandler, interceptors)
	if err != nil {
		return fmt.Errorf("failed to register ghcp service: %w", err)
	}

	log.Debugf("starting server on %s", serverAddress)
	err = http.ListenAndServe(serverAddress, h2c.NewHandler(router.Handler(), &http2.Server{}))
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func registerService(router dryhttp.Router, h api.Handler, opts ...connect.HandlerOption) error {
	path, handler, err := h.Build(opts...)
	if err != nil {
		return fmt.Errorf("failed to build handler: %s: %w", h.Name(), err)
	}

	router.AddRoute(dryhttp.ANY, fmt.Sprintf("%s*", path), handler)
	return nil
}

func buildConnectInterceptors() (connect.HandlerOption, error) {
	var interceptors []connect.Interceptor

	authInterceptor, err := api.NewAuthenticationInterceptor(api.AuthenticationInterceptorConfig{
		MockAuthentication: serverMockAuthentication,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create authentication interceptor: %w", err)
	}

	interceptors = append(interceptors, authInterceptor)

	validatorInterceptor, err := api.NewValidatorInterceptor()
	if err != nil {
		return nil, fmt.Errorf("failed to create validator interceptor: %w", err)
	}

	interceptors = append(interceptors, validatorInterceptor)

	return connect.WithInterceptors(interceptors...), nil
}
