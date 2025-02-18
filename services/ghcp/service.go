package ghcp

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	ghcpv1 "buf.build/gen/go/safedep/api/protocolbuffers/go/safedep/services/ghcp/v1"
	"github.com/safedep/dry/log"
	"github.com/safedep/ghcp/pkg/adapters/github"
	"github.com/safedep/ghcp/pkg/gh"
	"github.com/safedep/ghcp/services"
)

type GitHubCommentProxyServiceConfig struct {
	// If true, the service will only comment on public repositories.
	AllowOnlyPublicRepositories bool

	// If true, the service will only update comments that were created by the bot user.
	// This is to prevent misuse of the service to update comments created by other users.
	AllowOnlyOwnCommentUpdates bool

	// The username of the bot user configured to comment on the repository.
	// This is usually the username associated with the GITHUB_TOKEN
	BotUsername string
}

// Secure defaults for the GitHubCommentProxyServiceConfig
func DefaultGitHubCommentProxyServiceConfig() GitHubCommentProxyServiceConfig {
	return GitHubCommentProxyServiceConfig{
		AllowOnlyPublicRepositories: true,
		AllowOnlyOwnCommentUpdates:  true,
		BotUsername:                 "safedep-bot",
	}
}

type gitHubCommentProxyService struct {
	config         GitHubCommentProxyServiceConfig
	ghIssueAdapter github.GitHubIssueAdapter
}

var _ services.Service[*ghcpv1.CreatePullRequestCommentRequest,
	*ghcpv1.CreatePullRequestCommentResponse] = &gitHubCommentProxyService{}

func NewGitHubCommentProxyService(config GitHubCommentProxyServiceConfig,
	ghIssueAdapter github.GitHubIssueAdapter) (*gitHubCommentProxyService, error) {

	if config.AllowOnlyOwnCommentUpdates && config.BotUsername == "" {
		return nil, fmt.Errorf("bot username is required when AllowOnlyOwnCommentUpdates is true")
	}

	return &gitHubCommentProxyService{config: config, ghIssueAdapter: ghIssueAdapter}, nil
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

	prNumber, err := strconv.Atoi(request.GetPrNumber())
	if err != nil {
		return nil, fmt.Errorf("failed to convert pr number to int: %w", err)
	}

	if request.GetTag() == "" {
		return s.createNewComment(ctx, prNumber, request)
	}

	return s.updateExistingComment(ctx, prNumber, request)
}

func (s *gitHubCommentProxyService) createNewComment(ctx context.Context, prNumber int,
	request *ghcpv1.CreatePullRequestCommentRequest) (*ghcpv1.CreatePullRequestCommentResponse, error) {

	log.Debugf("Creating comment on PR: %s", request.GetPrNumber())
	comment, err := s.ghIssueAdapter.CreateIssueComment(ctx, request.GetOwner(),
		request.GetRepo(), prNumber, request.GetBody())
	if err != nil {
		return nil, fmt.Errorf("failed to create issue comment: %w", err)
	}

	return &ghcpv1.CreatePullRequestCommentResponse{
		CommentId: fmt.Sprintf("%d", comment.GetID()),
	}, nil
}

func (s *gitHubCommentProxyService) updateExistingComment(ctx context.Context, prNumber int,
	request *ghcpv1.CreatePullRequestCommentRequest) (*ghcpv1.CreatePullRequestCommentResponse, error) {

	log.Debugf("Updating comment on PR: %s with Tag: %s", request.GetPrNumber(), request.GetTag())

	comments, err := s.ghIssueAdapter.ListIssueComments(ctx, request.GetOwner(),
		request.GetRepo(), prNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to list issue comments: %w", err)
	}

	for _, comment := range comments {
		if strings.Contains(comment.GetBody(), request.GetTag()) {
			if s.config.AllowOnlyOwnCommentUpdates {
				if comment.GetUser().GetLogin() != s.config.BotUsername {
					return nil, fmt.Errorf("refusing to update comment created by another user")
				}
			}

			log.Debugf("Updating commentId: %d", comment.GetID())

			updatedComment, err := s.ghIssueAdapter.UpdateIssueComment(ctx, request.GetOwner(),
				request.GetRepo(), int(comment.GetID()), request.GetBody())
			if err != nil {
				return nil, fmt.Errorf("failed to update issue comment: %w", err)
			}

			return &ghcpv1.CreatePullRequestCommentResponse{
				CommentId: fmt.Sprintf("%d", updatedComment.GetID()),
			}, nil
		}
	}

	log.Debugf("No comment found with Tag: %s", request.GetTag())
	return nil, fmt.Errorf("no comment found with Tag: %s", request.GetTag())
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
