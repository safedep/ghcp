package ghcp

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	ghcpv1 "buf.build/gen/go/safedep/api/protocolbuffers/go/safedep/services/ghcp/v1"
	ghapi "github.com/google/go-github/v69/github"
	"github.com/safedep/ghcp/pkg/adapters/github"
	"github.com/safedep/ghcp/pkg/gh"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"
)

func TestGitHubCommentProxyService(t *testing.T) {
	cases := []struct {
		name             string
		config           GitHubCommentProxyServiceConfig
		token            *gh.GitHubTokenContext
		serviceInitError error
		mock             func(*github.MockGitHubIssueAdapter, *github.MockGitHubRepositoryAdapter)
		request          *ghcpv1.CreatePullRequestCommentRequest
		assert           func(*testing.T, error, *ghcpv1.CreatePullRequestCommentResponse)
	}{
		{
			name: "create new comment",
			config: GitHubCommentProxyServiceConfig{
				GitHubTokenAudienceName: GitHubTokenAudienceName,
			},
			token: &gh.GitHubTokenContext{
				Repository:      "safedep/ghcp",
				RepositoryOwner: "safedep",
				Audience:        GitHubTokenAudienceName,
			},
			mock: func(m *github.MockGitHubIssueAdapter, m2 *github.MockGitHubRepositoryAdapter) {
				m.EXPECT().CreateIssueComment(mock.Anything, "safedep", "ghcp", 1,
					"test comment").Return(&ghapi.IssueComment{ID: proto.Int64(1)}, nil)
			},
			request: &ghcpv1.CreatePullRequestCommentRequest{
				Owner:    "safedep",
				Repo:     "ghcp",
				PrNumber: "1",
				Body:     "test comment",
			},
			assert: func(t *testing.T, err error, res *ghcpv1.CreatePullRequestCommentResponse) {
				assert.NoError(t, err)
				assert.NotEmpty(t, res.GetCommentId())
			},
		},
		{
			name:   "create comment fails when no token is provided",
			config: GitHubCommentProxyServiceConfig{},
			request: &ghcpv1.CreatePullRequestCommentRequest{
				Owner:    "safedep",
				Repo:     "ghcp",
				PrNumber: "1",
			},
			assert: func(t *testing.T, err error, res *ghcpv1.CreatePullRequestCommentResponse) {
				assert.Error(t, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "create comment fails when audience in request does not match token context",
			config: GitHubCommentProxyServiceConfig{
				GitHubTokenAudienceName: GitHubTokenAudienceName,
			},
			token: &gh.GitHubTokenContext{
				Repository:      "safedep/ghcp",
				RepositoryOwner: "safedep",
				Audience:        "safedep-ghcp-test",
			},
			request: &ghcpv1.CreatePullRequestCommentRequest{
				Owner:    "safedep",
				Repo:     "ghcp",
				PrNumber: "1",
			},
			assert: func(t *testing.T, err error, res *ghcpv1.CreatePullRequestCommentResponse) {
				assert.Error(t, err)
				assert.Nil(t, res)
			},
		},
		{
			name:   "create comment fails when owner in request does not match token context",
			config: GitHubCommentProxyServiceConfig{},
			token: &gh.GitHubTokenContext{
				Repository:      "safedep/ghcp",
				RepositoryOwner: "safedep",
				Audience:        GitHubTokenAudienceName,
			},
			request: &ghcpv1.CreatePullRequestCommentRequest{
				Owner:    "safedep-test",
				Repo:     "ghcp",
				PrNumber: "1",
			},
			assert: func(t *testing.T, err error, res *ghcpv1.CreatePullRequestCommentResponse) {
				assert.Error(t, err)
				assert.Nil(t, res)
			},
		},
		{
			name:   "create comment fails when repo in request does not match token context",
			config: GitHubCommentProxyServiceConfig{},
			token: &gh.GitHubTokenContext{
				Repository:      "safedep/ghcp-test",
				RepositoryOwner: "safedep",
				Audience:        GitHubTokenAudienceName,
			},
			request: &ghcpv1.CreatePullRequestCommentRequest{
				Owner:    "safedep",
				Repo:     "ghcp",
				PrNumber: "1",
			},
			assert: func(t *testing.T, err error, res *ghcpv1.CreatePullRequestCommentResponse) {
				assert.Error(t, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "create comment fails when repo is is required to be public",
			config: GitHubCommentProxyServiceConfig{
				AllowOnlyPublicRepositories: true,
			},
			token: &gh.GitHubTokenContext{
				Repository:           "safedep/ghcp",
				RepositoryOwner:      "safedep",
				RepositoryVisibility: "private",
				Audience:             GitHubTokenAudienceName,
			},
			request: &ghcpv1.CreatePullRequestCommentRequest{
				Owner: "safedep",
				Repo:  "ghcp",
			},
			assert: func(t *testing.T, err error, res *ghcpv1.CreatePullRequestCommentResponse) {
				assert.Error(t, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "update comment is successful when tag is provided and comment exists",
			config: GitHubCommentProxyServiceConfig{
				GitHubTokenAudienceName: GitHubTokenAudienceName,
			},
			token: &gh.GitHubTokenContext{
				Repository:      "safedep/ghcp",
				RepositoryOwner: "safedep",
				Audience:        GitHubTokenAudienceName,
			},
			mock: func(m *github.MockGitHubIssueAdapter, m2 *github.MockGitHubRepositoryAdapter) {
				m.EXPECT().ListIssueComments(mock.Anything, "safedep", "ghcp", 1).
					Return([]*ghapi.IssueComment{
						{ID: proto.Int64(1), Body: proto.String("test comment with tag: test-tag")},
					}, nil)
				m.EXPECT().UpdateIssueComment(mock.Anything, "safedep", "ghcp", 1,
					"test comment").Return(&ghapi.IssueComment{ID: proto.Int64(1)}, nil)
			},
			request: &ghcpv1.CreatePullRequestCommentRequest{
				Owner:    "safedep",
				Repo:     "ghcp",
				PrNumber: "1",
				Body:     "test comment",
				Tag:      "test-tag",
			},
			assert: func(t *testing.T, err error, res *ghcpv1.CreatePullRequestCommentResponse) {
				assert.NoError(t, err)
				assert.NotEmpty(t, res.GetCommentId())
			},
		},
		{
			name: "update comment fails when when no comment found with tag",
			config: GitHubCommentProxyServiceConfig{
				GitHubTokenAudienceName: GitHubTokenAudienceName,
			},
			token: &gh.GitHubTokenContext{
				Repository:      "safedep/ghcp",
				RepositoryOwner: "safedep",
				Audience:        GitHubTokenAudienceName,
			},
			request: &ghcpv1.CreatePullRequestCommentRequest{
				Owner:    "safedep",
				Repo:     "ghcp",
				PrNumber: "1",
				Body:     "test comment",
				Tag:      "test-tag",
			},
			mock: func(m *github.MockGitHubIssueAdapter, m2 *github.MockGitHubRepositoryAdapter) {
				m.EXPECT().ListIssueComments(mock.Anything, "safedep", "ghcp", 1).
					Return([]*ghapi.IssueComment{}, nil)
			},
			assert: func(t *testing.T, err error, res *ghcpv1.CreatePullRequestCommentResponse) {
				assert.Error(t, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "update comment fails when the comment user is not same as the bot user",
			config: GitHubCommentProxyServiceConfig{
				AllowOnlyOwnCommentUpdates: true,
				BotUsername:                "safedep-bot",
				GitHubTokenAudienceName:    GitHubTokenAudienceName,
			},
			token: &gh.GitHubTokenContext{
				Repository:      "safedep/ghcp",
				RepositoryOwner: "safedep",
				Audience:        GitHubTokenAudienceName,
			},
			mock: func(m *github.MockGitHubIssueAdapter, m2 *github.MockGitHubRepositoryAdapter) {
				m.EXPECT().ListIssueComments(mock.Anything, "safedep", "ghcp", 1).
					Return([]*ghapi.IssueComment{
						{
							ID:   proto.Int64(1),
							User: &ghapi.User{Login: proto.String("test-user")},
							Body: proto.String("test comment with tag: test-tag"),
						},
					}, nil)

				// We do not expect any calls to UpdateIssueComment because the user is not the same
				// as the bot user. The service is expected to fail fast.
			},
			request: &ghcpv1.CreatePullRequestCommentRequest{
				Owner:    "safedep",
				Repo:     "ghcp",
				PrNumber: "1",
				Body:     "test comment",
				Tag:      "test-tag",
			},
			assert: func(t *testing.T, err error, res *ghcpv1.CreatePullRequestCommentResponse) {
				assert.Error(t, err)
				assert.Nil(t, res)
			},
		},
		{
			name:             "create comment fails when both installation verifiers and workload identity token verification are skipped",
			serviceInitError: fmt.Errorf("skip workload identity token verification is true but no installation verifiers are provided"),
			config: GitHubCommentProxyServiceConfig{
				SkipWorkloadIdentityTokenVerification: true,
			},
			request: &ghcpv1.CreatePullRequestCommentRequest{
				Owner:    "safedep",
				Repo:     "ghcp",
				PrNumber: "1",
				Body:     "test comment",
			},
			assert: func(t *testing.T, err error, res *ghcpv1.CreatePullRequestCommentResponse) {
				assert.Error(t, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "create comment is successful when workload identity token verification is skipped and installation verifiers are provided",
			config: GitHubCommentProxyServiceConfig{
				SkipWorkloadIdentityTokenVerification: true,
				InstallationVerifiers: []GitHubCommentsProxyInstallationVerifier{
					{
						Path:   "/.github/workflows/test-path",
						Action: regexp.MustCompile("test-content"),
					},
				},
			},
			request: &ghcpv1.CreatePullRequestCommentRequest{
				Owner:    "safedep",
				Repo:     "ghcp",
				PrNumber: "1",
				Body:     "test comment",
			},
			mock: func(m *github.MockGitHubIssueAdapter, m2 *github.MockGitHubRepositoryAdapter) {
				m2.EXPECT().GetFileContent(mock.Anything, "safedep", "ghcp", "/.github/workflows/test-path").
					Return([]byte("test-content"), nil)
				m.EXPECT().CreateIssueComment(mock.Anything, "safedep", "ghcp", 1,
					"test comment").Return(&ghapi.IssueComment{ID: proto.Int64(1)}, nil)
			},
			assert: func(t *testing.T, err error, res *ghcpv1.CreatePullRequestCommentResponse) {
				assert.NoError(t, err)
				assert.NotEmpty(t, res.GetCommentId())
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ghIssueAdapter := github.NewMockGitHubIssueAdapter(t)
			ghRepoAdapter := github.NewMockGitHubRepositoryAdapter(t)

			service, err := NewGitHubCommentProxyService(c.config, ghIssueAdapter, ghRepoAdapter)
			if c.serviceInitError != nil {
				assert.Error(t, err)
				assert.Nil(t, service)
				assert.Contains(t, err.Error(), c.serviceInitError.Error())
				return
			}

			assert.NoError(t, err)

			ctx := context.Background()
			if c.token != nil {
				ctx = gh.InjectGitHubTokenContext(ctx, *c.token)
			}

			if c.mock != nil {
				c.mock(ghIssueAdapter, ghRepoAdapter)
			}

			response, err := service.Execute(ctx, c.request)
			c.assert(t, err, response)
		})
	}
}
