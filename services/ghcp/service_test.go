package ghcp

import (
	"context"
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
		name    string
		config  GitHubCommentProxyServiceConfig
		token   *gh.GitHubTokenContext
		mock    func(*github.MockGitHubIssueAdapter)
		request *ghcpv1.CreatePullRequestCommentRequest
		assert  func(*testing.T, error, *ghcpv1.CreatePullRequestCommentResponse)
	}{
		{
			name:   "create new comment",
			config: GitHubCommentProxyServiceConfig{},
			token: &gh.GitHubTokenContext{
				Repository:      "safedep/ghcp",
				RepositoryOwner: "safedep",
			},
			mock: func(m *github.MockGitHubIssueAdapter) {
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
			name: "create comment fails when repo is is required to be public",
			config: GitHubCommentProxyServiceConfig{
				AllowOnlyPublicRepositories: true,
			},
			token: &gh.GitHubTokenContext{
				Repository:           "safedep/ghcp",
				RepositoryOwner:      "safedep",
				RepositoryVisibility: "private",
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
			name:   "update comment is successful when tag is provided and comment exists",
			config: GitHubCommentProxyServiceConfig{},
			token: &gh.GitHubTokenContext{
				Repository:      "safedep/ghcp",
				RepositoryOwner: "safedep",
			},
			mock: func(m *github.MockGitHubIssueAdapter) {
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
			name:   "update comment fails when when no comment found with tag",
			config: GitHubCommentProxyServiceConfig{},
			token: &gh.GitHubTokenContext{
				Repository:      "safedep/ghcp",
				RepositoryOwner: "safedep",
			},
			request: &ghcpv1.CreatePullRequestCommentRequest{
				Owner:    "safedep",
				Repo:     "ghcp",
				PrNumber: "1",
				Body:     "test comment",
				Tag:      "test-tag",
			},
			mock: func(m *github.MockGitHubIssueAdapter) {
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
			},
			token: &gh.GitHubTokenContext{
				Repository:      "safedep/ghcp",
				RepositoryOwner: "safedep",
			},
			mock: func(m *github.MockGitHubIssueAdapter) {
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
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ghIssueAdapter := github.NewMockGitHubIssueAdapter(t)
			service, err := NewGitHubCommentProxyService(c.config, ghIssueAdapter)
			assert.NoError(t, err)

			ctx := context.Background()
			if c.token != nil {
				ctx = gh.InjectGitHubTokenContext(ctx, *c.token)
			}

			if c.mock != nil {
				c.mock(ghIssueAdapter)
			}

			response, err := service.Execute(ctx, c.request)
			c.assert(t, err, response)
		})
	}
}
