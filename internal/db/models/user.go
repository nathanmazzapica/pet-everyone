package models

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
	ID        uuid.UUID `sql:"user_id"`
	Email     string    `sql:"email"`
	CreatedAt time.Time `sql:"created_at"`
}

func (u *UserModel) Create(email, password string) (*User, error) {

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
		ID:        uuid.New(),
		Email:     email,
		CreatedAt: time.Now(),
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	query := `INSERT INTO RegisteredUser (user_id, email, password_hash) VALUES ($1, $2, $3)`
	_, err = u.DB.Exec(query, user.ID, user.Email, hash)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
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
	query := `SELECT user_id, email, created_at FROM RegisteredUser WHERE user_id = $1`
	row := u.DB.QueryRow(query, id)

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.CreatedAt)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	return &user, err
}

func (u *UserModel) GetByEmail(email string) (*User, error) {
	query := `SELECT user_id, email, created_at FROM RegisteredUser WHERE email = $1`
	row := u.DB.QueryRow(query, email)

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.CreatedAt)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	return &user, err
}

func GetPetCountByPetID(db *sql.DB, id *uuid.UUID, petID *string) (int, error) {
	query := `SELECT COALESCE(SUM(click_count), 0) FROM UserPetsClickCount WHERE pet_id = $1 AND user_id = $2;`
	row := db.QueryRow(query, petID, id)

	var count int
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
