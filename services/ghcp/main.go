package ghcp

import (
	"context"
	"fmt"
	"strings"

	ghcpv1 "buf.build/gen/go/safedep/api/protocolbuffers/go/safedep/services/ghcp/v1"
	"github.com/safedep/ghcp/pkg/gh"
	"github.com/safedep/ghcp/services"
)

type GitHubCommentProxyServiceConfig struct {
	GitHubToken                 string
	AllowOnlyPublicRepositories bool
}

type gitHubCommentProxyService struct {
	config GitHubCommentProxyServiceConfig
}

var _ services.Service[*ghcpv1.CreatePullRequestCommentRequest,
	*ghcpv1.CreatePullRequestCommentResponse] = &gitHubCommentProxyService{}

func NewGitHubCommentProxyService(config GitHubCommentProxyServiceConfig) (*gitHubCommentProxyService, error) {
	return &gitHubCommentProxyService{config: config}, nil
}

func (s *gitHubCommentProxyService) Name() string {
	return "GitHubCommentProxyService"
}

func (s *gitHubCommentProxyService) Config() services.ServiceConfiguration {
	return services.ServiceConfiguration{}
}

func (s *gitHubCommentProxyService) Execute(ctx context.Context,
	request *ghcpv1.CreatePullRequestCommentRequest) (*ghcpv1.CreatePullRequestCommentResponse, error) {

	tokenContext, err := gh.ExtractGitHubTokenContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to extract GitHub Workload Identity Token context: %w", err)
	}

	if err := s.verifyRepositoryAccess(tokenContext, request); err != nil {
		return nil, fmt.Errorf("failed to verify repository access: %w", err)
	}

	return &ghcpv1.CreatePullRequestCommentResponse{}, nil
}

// verifyRepositoryAccess verifies that the token context matches the requested repository
// This is to prevent the service from being misused to spam comments to various repositories
func (s *gitHubCommentProxyService) verifyRepositoryAccess(tokenContext gh.GitHubTokenContext,
	req *ghcpv1.CreatePullRequestCommentRequest) error {

	if !strings.EqualFold(tokenContext.RepositoryOwner, req.GetOwner()) {
		return fmt.Errorf("repository owner mismatch: %s != %s", tokenContext.RepositoryOwner, req.GetOwner())
	}

	expectedRepository := fmt.Sprintf("%s/%s", req.GetOwner(), req.GetRepo())
	if !strings.EqualFold(tokenContext.Repository, expectedRepository) {
		return fmt.Errorf("repository mismatch: %s != %s", tokenContext.Repository, expectedRepository)
	}

	if s.config.AllowOnlyPublicRepositories && tokenContext.RepositoryVisibility != "public" {
		return fmt.Errorf("repository is not public: %s", tokenContext.RepositoryVisibility)
	}

	return nil
}
