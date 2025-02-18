package github

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitHubClientAdapterListIssueComments(t *testing.T) {
	cases := []struct {
		name        string
		org         string
		repo        string
		number      int
		minComments int
		err         error
	}{
		{
			name:        "valid repo",
			org:         "safedep",
			repo:        "vet",
			number:      349,
			minComments: 1,
		},
		{
			name:   "invalid repo",
			org:    "safedep",
			repo:   "invalid",
			number: 1,
			err:    errors.New("Not Found"),
		},
	}

	client, err := NewGitHubAdapter(DefaultGitHubAdapterConfig())
	if err != nil {
		t.Fatalf("failed to create github client: %v", err)
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			res, err := client.ListIssueComments(context.Background(), test.org, test.repo, test.number)
			if test.err != nil {
				assert.Error(t, err)
				assert.ErrorContains(t, err, test.err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.GreaterOrEqual(t, len(res), test.minComments)
			}
		})
	}
}
