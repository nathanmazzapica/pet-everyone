package auth

import (
	"errors"
	"testing"
)

func TestAuth_PasswordHash(t *testing.T) {

	tests := []struct {
		name          string
		password      string
		expectedError error
	}{
		{name: "valid password test", password: "test1234567890", expectedError: nil},
		{name: "too short password test", password: "test123", expectedError: TooShortPasswordError},
		{name: "empty password test", password: "", expectedError: EmptyPasswordError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)

			if !errors.Is(err, tt.expectedError) {
				t.Fatalf("Expected error: %v, got: %v", tt.expectedError, err)
			}

			if tt.expectedError == nil && hash == "" {
				t.Fatalf("Expected hash to be non-empty, got empty string")
			}
		})
	}
}

func TestAuth_CheckPasswordHash(t *testing.T) {

	tests := []struct {
		name     string
		password string
	}{
		{name: "normal password test", password: "test1235678101112"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, _ := HashPassword(tt.password)

			match, err := CheckPasswordHash(tt.password, hash)
			if err != nil {
				t.Fatalf("Error checking password (%s) against hash: %s\n: %v", tt.password, hash, err)
			}

			if !match {
				t.Fatalf("Password (%s) did not match hash", tt.password)
			}
		})
	}
}
