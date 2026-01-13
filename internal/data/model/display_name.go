package model

import (
	"database/sql"
	"errors"
	"pet-everyone/internal/displayname"
	"strings"

	"github.com/google/uuid"
)

const (
	ownerTypeUser  = "user"
	ownerTypeGuest = "guest"
)

var (
	ErrDisplayNameTaken = errors.New("display name already in use")
	errUnknownOwnerType = errors.New("unknown display name owner type")
)

func insertDisplayName(tx *sql.Tx, name string, ownerType string, ownerID uuid.UUID) error {
	if err := validateOwnerType(ownerType); err != nil {
		return err
	}
	if err := displayname.Validate(name); err != nil {
		return err
	}
	_, err := tx.Exec(`INSERT INTO DisplayName (name, owner_type, owner_id) VALUES ($1, $2, $3)`, name, ownerType, ownerID)
	if err != nil {
		if isUniqueViolation(err) {
			return ErrDisplayNameTaken
		}
		return err
	}
	return nil
}

func upsertDisplayName(tx *sql.Tx, name string, ownerType string, ownerID uuid.UUID) error {
	if err := validateOwnerType(ownerType); err != nil {
		return err
	}
	if err := displayname.Validate(name); err != nil {
		return err
	}
	_, err := tx.Exec(`
INSERT INTO DisplayName (name, owner_type, owner_id)
VALUES ($1, $2, $3)
ON CONFLICT (owner_type, owner_id) DO UPDATE SET name = EXCLUDED.name
`, name, ownerType, ownerID)
	if err != nil {
		if isUniqueViolation(err) {
			return ErrDisplayNameTaken
		}
		return err
	}
	return nil
}

func validateOwnerType(ownerType string) error {
	switch ownerType {
	case ownerTypeUser, ownerTypeGuest:
		return nil
	default:
		return errUnknownOwnerType
	}
}

func isUniqueViolation(err error) bool {
	return strings.Contains(err.Error(), "duplicate key value")
}
