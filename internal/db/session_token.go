package db

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrSessionTokenNotFound = errors.New("session token not found")
)

type SessionToken struct {
	Token  string    `sql:"token"`
	UserID string    `sql:"user_id"`
	Expiry time.Time `sql:"expires_at"`
}

func (c *Client) SaveSessionToken(token string, userID string, expiry time.Time) error {
	query := `INSERT INTO SessionTokens (token, user_id, expires_at) VALUES ($1, $2, $3)`
	_, err := c.db.Exec(query, token, userID, expiry)
	return err
}

func (c *Client) DeleteSessionToken(token string) error {
	query := `DELETE FROM SessionTokens WHERE token = $1`
	_, err := c.db.Exec(query, token)
	return err
}

func (c *Client) GetSessionToken(token string) (*SessionToken, error) {
	query := `SELECT * FROM SessionTokens WHERE token = $1`
	var session SessionToken
	err := c.db.QueryRow(query, token).Scan(&session.Token, &session.UserID, &session.Expiry)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSessionTokenNotFound
		}
	}
	return &session, err
}
