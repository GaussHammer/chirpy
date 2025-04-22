package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestJWTCreationAndValidation(t *testing.T) {
	// Create a test user ID
	userID := uuid.New()
	tokenSecret := "test-secret-key"

	// Test case 1: Valid token
	t.Run("Valid token", func(t *testing.T) {
		// Create a token that lasts for 1 hour
		token, err := MakeJWT(userID, tokenSecret, time.Hour)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// Validate the token
		extractedID, err := ValidateJWT(token, tokenSecret)
		assert.NoError(t, err)
		assert.Equal(t, userID, extractedID)
	})

	// Test case 2: Invalid secret
	t.Run("Invalid secret", func(t *testing.T) {
		// Create a token with one secret
		token, err := MakeJWT(userID, tokenSecret, time.Hour)
		assert.NoError(t, err)

		// Try to validate with a different secret
		wrongSecret := "wrong-secret"
		_, err = ValidateJWT(token, wrongSecret)
		assert.Error(t, err)
	})
}
