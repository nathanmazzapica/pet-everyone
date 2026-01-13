package model

import (
	"database/sql"
	"errors"
	"log"
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
	Expiry time.Time `sql:"expires_at"`
	UserID string    `sql:"user_id"`
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
	err := s.DB.QueryRow(query, token).Scan(&session.Token, &session.Expiry, &session.UserID)
	if err != nil {
		log.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSessionTokenNotFound
		}
	}
	return &session, err
}

func (s *SessionTokenModel) GetUserID(token string) (string, error) {
	session, err := s.Get(token)
	if err != nil {
		return "", err
	}
	return session.UserID, nil
}
