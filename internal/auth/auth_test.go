package auth

import (
	"testing"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	secret := "my-super-secret-key"
	userID := uuid.New()

	// 1. Test Valid Token
	t.Run("Valid Token", func(t *testing.T) {
		token, err := MakeJWT(userID, secret)
		if err != nil {
			t.Fatalf("failed to make JWT: %v", err)
		}

		parsedID, err := ValidateJWT(token, secret)
		if err != nil {
			t.Fatalf("failed to validate valid JWT: %v", err)
		}

		if parsedID != userID {
			t.Errorf("expected ID %v, got %v", userID, parsedID)
		}
	})

	// 2. Test Expired Token
	t.Run("Expired Token", func(t *testing.T) {
		// Set expiration to a negative duration (already expired)
		token, err := MakeJWT(userID, secret)
		if err != nil {
			t.Fatalf("failed to make JWT: %v", err)
		}

		_, err = ValidateJWT(token, secret)
		if err == nil {
			t.Error("expected error for expired token, but got none")
		}
	})

	// 3. Test Wrong Secret
	t.Run("Wrong Secret", func(t *testing.T) {
		token, err := MakeJWT(userID, "correct-secret")
		if err != nil {
			t.Fatalf("failed to make JWT: %v", err)
		}

		_, err = ValidateJWT(token, "wrong-secret")
		if err == nil {
			t.Error("expected error for wrong secret, but got none")
		}
	})
}
