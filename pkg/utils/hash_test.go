package utils

import "testing"

func TestHashPassword(t *testing.T) {
	password := "mySecretPassword"

	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if len(hashedPassword) == 0 {
		t.Errorf("Hashed password should not be empty")
	}

	// Check if the hashed password can be verified against the original password
	if !CheckPasswordHash(password, hashedPassword) {
		t.Errorf("Password hash mismatch: the hash does not match the original password")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "mySecretPassword"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Correct password check
	if !CheckPasswordHash(password, hashedPassword) {
		t.Errorf("Password hash mismatch: the hash does not match the original password")
	}

	// Incorrect password check
	wrongPassword := "wrongPassword"
	if CheckPasswordHash(wrongPassword, hashedPassword) {
		t.Errorf("Password hash should not match the incorrect password")
	}
}
