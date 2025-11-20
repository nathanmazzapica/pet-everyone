package auth

import (
	"errors"
	"fmt"

	"github.com/alexedwards/argon2id"
)

const MinPasswordLength = 8

var (
	EmptyPasswordError    = errors.New("password cannot be empty")
	TooShortPasswordError = errors.New(fmt.Sprintf("password must be at least %d characters long", MinPasswordLength))
)

func HashPassword(password string) (string, error) {

	if len(password) == 0 {
		return "", EmptyPasswordError
	}

	if len(password) < MinPasswordLength {
		return "", TooShortPasswordError
	}

	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	if hash == "" {
	}
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return match, nil
}
