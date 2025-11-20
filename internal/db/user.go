package db

import (
	"database/sql"
	"errors"
	"fmt"
	"pet-everyone/internal/auth"
	"time"

	"github.com/google/uuid"
	"net/mail"
)

var (
	ErrUserNotFound = fmt.Errorf("user not found")
	ErrEmailInUse   = fmt.Errorf("email address already in use")
	ErrEmptyEmail   = fmt.Errorf("email cannot be empty")
	ErrInvalidEmail = fmt.Errorf("invalid email address")
)

type User struct {
	ID        uuid.UUID `sql:"user_id"`
	Email     string    `sql:"email"`
	CreatedAt time.Time `sql:"created_at"`
}

func (c *Client) CreateUser(email, password string) (*User, error) {

	if len(email) == 0 {
		return nil, ErrEmptyEmail
	}

	if !isValidEmail(email) {
		return nil, ErrInvalidEmail
	}

	exists, err := c.checkEmailExists(email)
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
	_, err = c.db.Exec(query, user.ID, user.Email, hash)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (c *Client) GetUserPasswordHashByEmail(email string) (string, error) {
	query := `SELECT password_hash FROM RegisteredUser WHERE email = $1`
	var hash string
	err := c.db.QueryRow(query, email).Scan(&hash)
	return hash, err
}

func (c *Client) GetUserByID(id *uuid.UUID) (*User, error) {
	query := `SELECT user_id, email, created_at FROM RegisteredUser WHERE user_id = $1`
	row := c.db.QueryRow(query, id)

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.CreatedAt)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	return &user, err
}

func (c *Client) GetUserByEmail(email string) (*User, error) {
	query := `SELECT user_id, email, created_at FROM RegisteredUser WHERE email = $1`
	row := c.db.QueryRow(query, email)

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.CreatedAt)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	return &user, err
}

func (c *Client) DeleteUser(id *uuid.UUID) error {
	query := `DELETE FROM RegisteredUser WHERE user_id = $1`
	_, err := c.db.Exec(query, id)
	return err
}

func (c *Client) checkEmailExists(email string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM RegisteredUser WHERE email = $1)`
	var exists bool
	err := c.db.QueryRow(query, email).Scan(&exists)
	return exists, err
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
