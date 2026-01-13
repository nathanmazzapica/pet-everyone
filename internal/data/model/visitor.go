package model

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrVisitorNotFound = errors.New("visitor not found")
)

type VisitorModel struct {
	DB *sql.DB
}

type Visitor struct {
	ID        uuid.UUID `sql:"guest_id"`
	CreatedAt time.Time `sql:"created_at"`
	LastSeen  time.Time `sql:"last_seen"`
}

func (m *VisitorModel) Create(id uuid.UUID) (*Visitor, error) {
	now := time.Now()
	query := `INSERT INTO Visitor (guest_id, created_at, last_seen) VALUES ($1, $2, $2) RETURNING guest_id, created_at, last_seen`

	var v Visitor
	err := m.DB.QueryRow(query, id, now).Scan(&v.ID, &v.CreatedAt, &v.LastSeen)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (m *VisitorModel) Get(id uuid.UUID) (*Visitor, error) {
	query := `SELECT guest_id, created_at, last_seen FROM Visitor WHERE guest_id = $1`
	row := m.DB.QueryRow(query, id)

	var v Visitor
	err := row.Scan(&v.ID, &v.CreatedAt, &v.LastSeen)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrVisitorNotFound
		}
		return nil, err
	}

	return &v, nil
}

func (m *VisitorModel) Upsert(id uuid.UUID) (*Visitor, error) {
	now := time.Now()
	query := `
INSERT INTO Visitor (guest_id, created_at, last_seen)
VALUES ($1, $2, $2)
ON CONFLICT (guest_id) DO UPDATE SET last_seen = EXCLUDED.last_seen
RETURNING guest_id, created_at, last_seen`

	var v Visitor
	err := m.DB.QueryRow(query, id, now).Scan(&v.ID, &v.CreatedAt, &v.LastSeen)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (m *VisitorModel) UpdateLastSeen(id uuid.UUID, at time.Time) error {
	query := `UPDATE Visitor SET last_seen = $2 WHERE guest_id = $1`
	_, err := m.DB.Exec(query, id, at)
	if err != nil {
		return err
	}
	return nil
}

func (m *VisitorModel) GetPetCountByPetID(guestID *uuid.UUID, petID *string) (int64, error) {
	query := `SELECT COALESCE(SUM(click_count), 0) FROM UserPetsClickCount WHERE pet_id = $1 AND guest_id = $2;`
	row := m.DB.QueryRow(query, petID, guestID)

	var count int64
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, err
}
