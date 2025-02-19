package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPAT(t *testing.T) {
	s := &authenticationInterceptor{}

	t.Run("should return true if the token starts with a valid PAT prefix", func(t *testing.T) {
		assert.True(t, s.isPAT("ghp_1234567890"))
		assert.True(t, s.isPAT("gho_1234567890"))
		assert.True(t, s.isPAT("ghu_1234567890"))
		assert.True(t, s.isPAT("ghs_1234567890"))
	})

	t.Run("should return false if the token does not start with a valid PAT prefix", func(t *testing.T) {
		assert.False(t, s.isPAT("1234567890"))
	})
}
