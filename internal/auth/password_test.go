package auth

import (
	"testing"
)

// TestHashAndCheckPassword ensures that HashPassword and CheckPasswordHash
// work correctly for valid and invalid passwords.
func TestHashAndCheckPassword(t *testing.T) {
	password := "mysecret"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if hash == "" {
		t.Fatal("HashPassword returned empty hash")
	}

	// Valid password should succeed
	if err := CheckPasswordHash(hash, password); err != nil {
		t.Errorf("CheckPasswordHash returned error for valid password: %v", err)
	}

	// Invalid password should fail
	wrong := "wrongpassword"
	if err := CheckPasswordHash(hash, wrong); err == nil {
		t.Error("CheckPasswordHash did not return error for invalid password")
	}
}

// TestHashPasswordUniqueSalt ensures that hashing the same password twice
// produces different hashes (due to salting) and both validate correctly.
func TestHashPasswordUniqueSalt(t *testing.T) {
	password := "anothersecret"
	hash1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Error hashing password first time: %v", err)
	}
	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Error hashing password second time: %v", err)
	}

	if hash1 == hash2 {
		t.Error("Expected different hashes for the same password, got identical values")
	}

	// Both hashes should validate correctly
	if err := CheckPasswordHash(hash1, password); err != nil {
		t.Errorf("CheckPasswordHash failed for hash1: %v", err)
	}
	if err := CheckPasswordHash(hash2, password); err != nil {
		t.Errorf("CheckPasswordHash failed for hash2: %v", err)
	}
}
