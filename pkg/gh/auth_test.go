package gh

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInjectGitHubTokenContext(t *testing.T) {
	tokenContext := GitHubTokenContext{
		Repository: "safedep/ghcp",
	}

	t.Run("should inject token context into context", func(t *testing.T) {
		ctx := context.Background()
		ctx = InjectGitHubTokenContext(ctx, tokenContext)

		extractedTokenContext, err := ExtractGitHubTokenContext(ctx)
		assert.NoError(t, err)
		assert.Equal(t, tokenContext, extractedTokenContext)
	})

	t.Run("should return error if no token context is found", func(t *testing.T) {
		ctx := context.Background()
		_, err := ExtractGitHubTokenContext(ctx)
		assert.Error(t, err)
	})
}
