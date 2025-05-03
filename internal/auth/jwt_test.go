package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

const testSecret = "test-secret"

func TestMakeJWT(t *testing.T) {
	// Generate a UUID for the user
	userID := uuid.New()

	// Set the expiration duration for the token
	expiresIn := time.Hour

	// Call MakeJWT to create a new token
	token, err := MakeJWT(userID, testSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	// Ensure the token is not empty
	if token == "" {
		t.Fatal("Expected a non-empty token")
	}
}

func TestValidateJWT(t *testing.T) {
	// Generate a UUID for the user
	userID := uuid.New()

	// Set the expiration duration for the token
	expiresIn := time.Hour

	// Create a valid token using MakeJWT
	token, err := MakeJWT(userID, testSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	// Validate the token with the correct secret
	validUserID, err := ValidateJWT(token, testSecret)
	if err != nil {
		t.Fatalf("ValidateJWT returned error: %v", err)
	}

	// Ensure the user ID matches
	if validUserID != userID {
		t.Fatalf("Expected user ID %v, got %v", userID, validUserID)
	}
}

func TestValidateExpiredJWT(t *testing.T) {
	// Generate a UUID for the user
	userID := uuid.New()

	// Set the expiration duration for the token to a past time
	expiresIn := -time.Hour // expired by 1 hour

	// Create an expired token using MakeJWT
	token, err := MakeJWT(userID, testSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	// Validate the expired token with the correct secret
	_, err = ValidateJWT(token, testSecret)
	if err == nil {
		t.Fatal("Expected error for expired token, but got nil")
	}
}

func TestValidateJWTWithWrongSecret(t *testing.T) {
	// Generate a UUID for the user
	userID := uuid.New()

	// Set the expiration duration for the token
	expiresIn := time.Hour

	// Create a valid token using MakeJWT
	token, err := MakeJWT(userID, testSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	// Validate the token with a wrong secret
	_, err = ValidateJWT(token, "wrong-secret")
	if err == nil {
		t.Fatal("Expected error for invalid secret, but got nil")
	}
}
