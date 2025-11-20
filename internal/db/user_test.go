package db

import (
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func TestUser_CreateTestUser(t *testing.T) {
	client := getTestDB()
	user, err := client.CreateUser("email@email.com", "test1234567890")
	defer cleanup(&user.ID)

	if err != nil {
		t.Fatalf("Error creating test user: %s", err)
	}
}

func TestUser_TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{"email@email.com", true},
		{"@email.com", false},
		{"email@", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			result := isValidEmail(tt.email)
			if result != tt.valid {
				t.Errorf("Expected %t, got %t", tt.valid, result)
			}
		})
	}
}

func getTestDB() Client {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	config := LoadConfig()
	client, err := Connect(config)
	if err != nil {
		panic(err)
	}

	return client
}

func cleanup(user_ids ...*uuid.UUID) {
	config := LoadConfig()
	client, err := Connect(config)
	if err != nil {
		panic(err)
	}

	for _, id := range user_ids {
		err := client.DeleteUser(id)
		if err != nil {
			panic(err)
		}
	}
}
