package model

import (
	"database/sql"
	"errors"
	"fmt"
	"net/mail"
	"pet-everyone/internal/auth"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound = fmt.Errorf("user not found")
	ErrEmailInUse   = fmt.Errorf("email address already in use")
	ErrEmptyEmail   = fmt.Errorf("email cannot be empty")
	ErrInvalidEmail = fmt.Errorf("invalid email address")
)

type UserModel struct {
	DB *sql.DB
}

type User struct {
	ID          uuid.UUID `sql:"user_id" json:"user_id"`
	Email       string    `sql:"email" json:"email"`
	DisplayName string    `sql:"display_name" json:"display_name"`
	CreatedAt   time.Time `sql:"created_at" json:"created_at"`
}

func (u *UserModel) Create(email, password, displayName string) (*User, error) {

	if len(email) == 0 {
		return nil, ErrEmptyEmail
	}

	if !isValidEmail(email) {
		return nil, ErrInvalidEmail
	}

	exists, err := u.checkEmailExists(email)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, ErrEmailInUse
	}

	user := &User{
		ID:          uuid.New(),
		Email:       email,
		DisplayName: displayName,
		CreatedAt:   time.Now(),
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	tx, err := u.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := insertDisplayName(tx, user.DisplayName, ownerTypeUser, user.ID); err != nil {
		return nil, err
	}

	query := `INSERT INTO RegisteredUser (user_id, email, password_hash, display_name) VALUES ($1, $2, $3, $4)`
	_, err = tx.Exec(query, user.ID, user.Email, hash, user.DisplayName)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserModel) GetPasswordHashByEmail(email string) (string, error) {
	query := `SELECT password_hash FROM RegisteredUser WHERE email = $1`
	var hash string
	err := u.DB.QueryRow(query, email).Scan(&hash)
	return hash, err
}

func (u *UserModel) Get(id *uuid.UUID) (*User, error) {
	query := `SELECT user_id, email, display_name, created_at FROM RegisteredUser WHERE user_id = $1`
	row := u.DB.QueryRow(query, id)

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.DisplayName, &user.CreatedAt)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	return &user, err
}

func (u *UserModel) GetByEmail(email string) (*User, error) {
	query := `SELECT user_id, email, display_name, created_at FROM RegisteredUser WHERE email = $1`
	row := u.DB.QueryRow(query, email)

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.DisplayName, &user.CreatedAt)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	return &user, err
}

func (u *UserModel) GetPetCountByPetID(userID *uuid.UUID, petID *string) (int64, error) {
	query := `SELECT COALESCE(SUM(click_count), 0) FROM UserPetsClickCount WHERE pet_id = $1 AND user_id = $2;`
	row := u.DB.QueryRow(query, petID, userID)

	var count int64
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, err
}

func (u *UserModel) Delete(id *uuid.UUID) error {
	query := `DELETE FROM RegisteredUser WHERE user_id = $1`
	_, err := u.DB.Exec(query, id)
	return err
}

// SetDisplayName updates the display name for a user, enforcing global uniqueness.
func (u *UserModel) SetDisplayName(userID uuid.UUID, displayName string) (string, error) {
	tx, err := u.DB.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	if err := upsertDisplayName(tx, displayName, ownerTypeUser, userID); err != nil {
		return "", err
	}

	res, err := tx.Exec(`UPDATE RegisteredUser SET display_name = $1 WHERE user_id = $2`, displayName, userID)
	if err != nil {
		return "", err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return "", err
	}
	if affected == 0 {
		return "", ErrUserNotFound
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	return displayName, nil
}

func (u *UserModel) checkEmailExists(email string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM RegisteredUser WHERE email = $1)`
	var exists bool
	err := u.DB.QueryRow(query, email).Scan(&exists)
	return exists, err
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
