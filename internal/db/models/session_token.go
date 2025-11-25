package models

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrSessionTokenNotFound = errors.New("session token not found")
)

type SessionTokenModel struct {
	DB *sql.DB
}

type SessionToken struct {
	Token  string    `sql:"token"`
	UserID string    `sql:"user_id"`
	Expiry time.Time `sql:"expires_at"`
}

func (t *SessionToken) IsExpired() bool {
	return time.Now().After(t.Expiry)
}

func (s *SessionTokenModel) Save(token string, userID string, expiry time.Time) error {
	query := `INSERT INTO SessionTokens (token, user_id, expires_at) VALUES ($1, $2, $3)`
	_, err := s.DB.Exec(query, token, userID, expiry)
	return err
}

func (s *SessionTokenModel) Delete(token string) error {
	query := `DELETE FROM SessionTokens WHERE token = $1`
	_, err := s.DB.Exec(query, token)
	return err
}

func (s *SessionTokenModel) Get(token string) (*SessionToken, error) {
	query := `SELECT * FROM SessionTokens WHERE token = $1`
	var session SessionToken
	err := s.DB.QueryRow(query, token).Scan(&session.Token, &session.UserID, &session.Expiry)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSessionTokenNotFound
		}
	}
	return &session, err
}
