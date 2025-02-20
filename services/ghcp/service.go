package ghcp

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	ghcpv1 "buf.build/gen/go/safedep/api/protocolbuffers/go/safedep/services/ghcp/v1"
	"github.com/safedep/dry/log"
	"github.com/safedep/dry/obs"
	"github.com/safedep/ghcp/pkg/adapters/github"
	"github.com/safedep/ghcp/pkg/gh"
	"github.com/safedep/ghcp/services"
)

const (
	GitHubTokenAudienceName = "safedep-ghcp"
	BotUsername             = "safedep-bot"
)

var (
	createCommentMetric              = obs.NewCounter("ghcp_create_comment_total", "Total number of comments created")
	updateCommentMetric              = obs.NewCounter("ghcp_update_comment_total", "Total number of comments updated")
	verifyInstallationMetric         = obs.NewCounter("ghcp_verify_installation_total", "Total number of installations verified")
	verifyRepositoryAccessMetric     = obs.NewCounter("ghcp_verify_repository_access_total", "Total number of repository accesses verified")
	successfulServiceExecutionMetric = obs.NewCounter("ghcp_successful_service_execution_total", "Total number of successful service executions")
	failedServiceExecutionMetric     = obs.NewCounter("ghcp_failed_service_execution_total", "Total number of failed service executions")
)

// GitHubCommentsProxyInstallationVerifier is a struct that verifies the installation of
// a given GitHub Action by path and action name. This is required to authenticate a repository
// when workload identity token is not available
type GitHubCommentsProxyInstallationVerifier struct {
	Path   string
	Action *regexp.Regexp
}

type GitHubCommentProxyServiceConfig struct {
	// If true, the service will only comment on public repositories.
	AllowOnlyPublicRepositories bool

	// If true, the service will only update comments that were created by the bot user.
	// This is to prevent misuse of the service to update comments created by other users.
	AllowOnlyOwnCommentUpdates bool

	// The username of the bot user configured to comment on the repository.
	// This is usually the username associated with the GITHUB_TOKEN
	BotUsername string

	// Audience name to verify against the GitHub Workload Identity Token
	GitHubTokenAudienceName string

	// Verify installation by checking the existence of a file
	VerifyInstallation bool

	// Verify vet installation
	InstallationVerifiers []GitHubCommentsProxyInstallationVerifier

	// InsecureSkipAuthorization is used to skip authorization checks for testing purposes
	InsecureSkipAuthorization bool
}

// Secure defaults for the GitHubCommentProxyServiceConfig
func DefaultGitHubCommentProxyServiceConfig() GitHubCommentProxyServiceConfig {
	return GitHubCommentProxyServiceConfig{
		AllowOnlyPublicRepositories: true,
		AllowOnlyOwnCommentUpdates:  true,
		BotUsername:                 BotUsername,
		GitHubTokenAudienceName:     GitHubTokenAudienceName,
		InstallationVerifiers: []GitHubCommentsProxyInstallationVerifier{
			{
				Path:   ".github/workflows/vet.yml",
				Action: regexp.MustCompile(`uses:\s+safedep/vet-action`),
			},
			{
				Path:   ".github/workflows/vet-ci.yml",
				Action: regexp.MustCompile(`uses:\s+safedep/vet-action`),
			},
		},
	}
}

type gitHubCommentProxyService struct {
	config         GitHubCommentProxyServiceConfig
	ghIssueAdapter github.GitHubIssueAdapter
	ghRepoAdapter  github.GitHubRepositoryAdapter
}

var _ services.Service[*ghcpv1.CreatePullRequestCommentRequest,
	*ghcpv1.CreatePullRequestCommentResponse] = &gitHubCommentProxyService{}

func NewGitHubCommentProxyService(config GitHubCommentProxyServiceConfig,
	ghIssueAdapter github.GitHubIssueAdapter,
	ghRepoAdapter github.GitHubRepositoryAdapter) (*gitHubCommentProxyService, error) {

	if config.AllowOnlyOwnCommentUpdates && config.BotUsername == "" {
		return nil, fmt.Errorf("bot username is required when AllowOnlyOwnCommentUpdates is true")
	}

	return &gitHubCommentProxyService{
		config:         config,
		ghIssueAdapter: ghIssueAdapter,
		ghRepoAdapter:  ghRepoAdapter,
	}, nil
}

func (s *gitHubCommentProxyService) Name() string {
	return "GitHubCommentProxyService"
}

func (s *gitHubCommentProxyService) Config() services.ServiceConfiguration {
	return services.ServiceConfiguration{}
}

func (s *gitHubCommentProxyService) Execute(ctx context.Context,
	request *ghcpv1.CreatePullRequestCommentRequest) (*ghcpv1.CreatePullRequestCommentResponse, error) {

	r, err := func() (*ghcpv1.CreatePullRequestCommentResponse, error) {
		if !s.config.InsecureSkipAuthorization {
			tokenContext, err := gh.ExtractGitHubTokenContext(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to extract GitHub Workload Identity Token context: %w", err)
			}

			if err := s.verifyRepositoryAccess(ctx, tokenContext, request); err != nil {
				return nil, fmt.Errorf("failed to verify repository access: %w", err)
			}
		}

		if s.config.VerifyInstallation {
			if err := s.verifyInstallation(ctx, request.GetOwner(), request.GetRepo()); err != nil {
				return nil, fmt.Errorf("failed to verify installation: %w", err)
			}
		}

		prNumber, err := strconv.Atoi(request.GetPrNumber())
		if err != nil {
			return nil, fmt.Errorf("failed to convert pr number to int: %w", err)
		}

		if request.GetTag() == "" {
			return s.createNewComment(ctx, prNumber, request)
		}

		return s.updateExistingComment(ctx, prNumber, request)
	}()

	if err != nil {
		log.Errorf("failed to execute service: %s", err)
		failedServiceExecutionMetric.Inc()
		return nil, err
	}

	successfulServiceExecutionMetric.Inc()
	return r, nil
}

func (s *gitHubCommentProxyService) createNewComment(ctx context.Context, prNumber int,
	request *ghcpv1.CreatePullRequestCommentRequest) (*ghcpv1.CreatePullRequestCommentResponse, error) {

	createCommentMetric.Inc()
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

	updateCommentMetric.Inc()
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
func (s *gitHubCommentProxyService) verifyRepositoryAccess(ctx context.Context, tokenContext gh.GitHubTokenContext,
	req *ghcpv1.CreatePullRequestCommentRequest) error {
	verifyRepositoryAccessMetric.Inc()

	if tokenContext.IsWorkloadIdentityToken() {
		return s.verifyWorkloadIdentityToken(tokenContext, req)
	}

	if tokenContext.IsActionToken() {
		return s.verifyActionToken(ctx, tokenContext, req)
	}

	return fmt.Errorf("failed to verify repository access for token context")
}

func (s *gitHubCommentProxyService) verifyActionToken(ctx context.Context, _ gh.GitHubTokenContext,
	req *ghcpv1.CreatePullRequestCommentRequest) error {
	prNumber, err := strconv.Atoi(req.GetPrNumber())
	if err != nil {
		return fmt.Errorf("failed to convert pr number to int: %w", err)
	}

	repo, err := s.ghRepoAdapter.GetRepository(ctx, req.GetOwner(), req.GetRepo())
	if err != nil {
		return fmt.Errorf("failed to get repository: %w", err)
	}

	if s.config.AllowOnlyPublicRepositories && repo.GetVisibility() != "public" {
		return fmt.Errorf("repository is not public")
	}

	pr, err := s.ghRepoAdapter.GetPullRequest(ctx, req.GetOwner(), req.GetRepo(), prNumber)
	if err != nil {
		return fmt.Errorf("failed to get pull request: %w", err)
	}

	if pr.GetState() != "open" {
		return fmt.Errorf("pull request is not open: %s", pr.GetState())
	}

	return nil
}

func (s *gitHubCommentProxyService) verifyWorkloadIdentityToken(tokenContext gh.GitHubTokenContext,
	req *ghcpv1.CreatePullRequestCommentRequest) error {
	if !strings.EqualFold(tokenContext.Audience, s.config.GitHubTokenAudienceName) {
		return fmt.Errorf("audience mismatch: %s != %s", tokenContext.Audience, s.config.GitHubTokenAudienceName)
	}

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

func (s *gitHubCommentProxyService) verifyInstallation(ctx context.Context, owner, repo string) error {
	verifyInstallationMetric.Inc()

	for _, verifier := range s.config.InstallationVerifiers {
		content, err := s.ghRepoAdapter.GetFileContent(ctx, owner, repo, verifier.Path)
		if err != nil {
			log.Debugf("verifyInstallation: %s/%s: failed to get file content: %s", owner, repo, err)
			continue
		}

		if verifier.Action.Match(content) {
			return nil
		}
	}

	return fmt.Errorf("no installation verifier matched")
}
