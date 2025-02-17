package jwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenGenerationAndValidation(t *testing.T) {
	username := "testuser"

	token, err := GenerateToken(username)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	validatedUsername, err := ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, username, validatedUsername)
}
