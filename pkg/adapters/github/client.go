package github

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/google/go-github/v69/github"
	"github.com/safedep/dry/log"
	"google.golang.org/protobuf/proto"
)

type GitHubAdapterConfig struct {
	// PAT / Token based authentication
	Token string

	// ClientId and ClientSecret for basic authentication
	// https://docs.github.com/en/rest/authentication/authenticating-to-the-rest-api#using-basic-authentication
	// App credentials usually have higher rate limits
	ClientId     string
	ClientSecret string

	// This is useful when we want to supply a client that
	// can handle rate limiting, etc.
	HTTPClient *http.Client
}

func DefaultGitHubAdapterConfig() GitHubAdapterConfig {
	token := os.Getenv("GHCP_GITHUB_TOKEN")
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}

	clientId, clientSecret := os.Getenv("GITHUB_CLIENT_ID"),
		os.Getenv("GITHUB_CLIENT_SECRET")

	return GitHubAdapterConfig{
		Token:        token,
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}
}

//go:generate mockery --name=GitHubIssueAdapter
type GitHubIssueAdapter interface {
	ListIssueComments(ctx context.Context, owner, repo string, number int) ([]*github.IssueComment, error)
	CreateIssueComment(ctx context.Context, owner, repo string, number int, comment string) (*github.IssueComment, error)
	UpdateIssueComment(ctx context.Context, owner, repo string, commentId int, comment string) (*github.IssueComment, error)
}

//go:generate mockery --name=GitHubRepositoryAdapter
type GitHubRepositoryAdapter interface {
	GetFileContent(ctx context.Context, owner, repo, path string) ([]byte, error)
}

type githubClient struct {
	client *github.Client
	config GitHubAdapterConfig
}

var _ GitHubIssueAdapter = &githubClient{}

type basicAuthTransportWrapper struct {
	Transport http.RoundTripper
	Username  string
	Password  string
}

func (b *basicAuthTransportWrapper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(b.Username, b.Password)
	return b.Transport.RoundTrip(req)
}

func NewGitHubAdapter(config GitHubAdapterConfig) (*githubClient, error) {
	if config.HTTPClient == nil {
		config.HTTPClient = http.DefaultClient
	}

	client := github.NewClient(config.HTTPClient)

	// Client credentials have highest precedence
	// for client authentication
	if config.ClientId != "" && config.ClientSecret != "" {
		log.Debugf("Using client credentials for GitHub authentication")
		client.Client().Transport = &basicAuthTransportWrapper{
			Transport: client.Client().Transport,
			Username:  config.ClientId,
			Password:  config.ClientSecret,
		}
	} else if config.Token != "" {
		log.Debugf("Using token for GitHub authentication")
		client = client.WithAuthToken(config.Token)
	} else {
		log.Warnf("Created a GitHub client without a token. This may cause rate limiting issues.")
	}

	return &githubClient{
		client: client,
		config: config,
	}, nil
}

func (g *githubClient) ListIssueComments(ctx context.Context, owner, repo string, number int) ([]*github.IssueComment, error) {
	comments, _, err := g.client.Issues.ListComments(ctx, owner, repo, int(number), &github.IssueListCommentsOptions{
		Sort:      proto.String("updated"),
		Direction: proto.String("desc"),
	})

	return comments, err
}

func (g *githubClient) CreateIssueComment(ctx context.Context, owner, repo string, number int, comment string) (*github.IssueComment, error) {
	issueComment, _, err := g.client.Issues.CreateComment(ctx, owner, repo, number, &github.IssueComment{Body: &comment})
	return issueComment, err
}

func (g *githubClient) UpdateIssueComment(ctx context.Context, owner, repo string, commentId int, comment string) (*github.IssueComment, error) {
	issueComment, _, err := g.client.Issues.EditComment(ctx, owner, repo, int64(commentId), &github.IssueComment{Body: &comment})
	return issueComment, err
}

func (g *githubClient) GetFileContent(ctx context.Context, owner, repo, path string) ([]byte, error) {
	content, _, _, err := g.client.Repositories.GetContents(ctx, owner, repo, path, &github.RepositoryContentGetOptions{})
	if err != nil {
		return nil, err
	}

	data, err := content.GetContent()
	if err != nil {
		return nil, fmt.Errorf("failed to get content of file %s: %w", path, err)
	}

	return []byte(data), nil
}
