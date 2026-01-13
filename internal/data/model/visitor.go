package model

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"pet-everyone/internal/displayname"

	"github.com/google/uuid"
)

var (
	ErrVisitorNotFound = errors.New("visitor not found")
)

type VisitorModel struct {
	DB *sql.DB
}

type Visitor struct {
	ID          uuid.UUID `sql:"guest_id"`
	DisplayName string    `sql:"display_name"`
	CreatedAt   time.Time `sql:"created_at"`
	LastSeen    time.Time `sql:"last_seen"`
}

func (m *VisitorModel) Create(id uuid.UUID) (*Visitor, error) {
	now := time.Now()
	name := displayname.Generate(nil)

	tx, err := m.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := insertDisplayName(tx, name, ownerTypeGuest, id); err != nil {
		return nil, err
	}

	query := `INSERT INTO Visitor (guest_id, created_at, last_seen, display_name) VALUES ($1, $2, $2, $3) RETURNING guest_id, display_name, created_at, last_seen`

	var v Visitor
	err = tx.QueryRow(query, id, now, name).Scan(&v.ID, &v.DisplayName, &v.CreatedAt, &v.LastSeen)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &v, nil
}

func (m *VisitorModel) Get(id uuid.UUID) (*Visitor, error) {
	query := `SELECT guest_id, display_name, created_at, last_seen FROM Visitor WHERE guest_id = $1`
	row := m.DB.QueryRow(query, id)

	var v Visitor
	err := row.Scan(&v.ID, &v.DisplayName, &v.CreatedAt, &v.LastSeen)
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

	tx, err := m.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var existing Visitor
	err = tx.QueryRow(`SELECT guest_id, display_name, created_at, last_seen FROM Visitor WHERE guest_id = $1`, id).
		Scan(&existing.ID, &existing.DisplayName, &existing.CreatedAt, &existing.LastSeen)

	switch {
	case err == nil:
		_, err = tx.Exec(`UPDATE Visitor SET last_seen = $2 WHERE guest_id = $1`, id, now)
		if err != nil {
			return nil, err
		}
		existing.LastSeen = now
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return &existing, nil
	case errors.Is(err, sql.ErrNoRows):
		var name string
		for attempts := 0; attempts < 5; attempts++ {
			name = displayname.Generate(nil)
			if err := insertDisplayName(tx, name, ownerTypeGuest, id); err != nil {
				if errors.Is(err, ErrDisplayNameTaken) {
					continue
				}
				return nil, err
			}
			break
		}
		if name == "" {
			return nil, fmt.Errorf("unable to generate unique guest display name")
		}

		var v Visitor
		err = tx.QueryRow(`
INSERT INTO Visitor (guest_id, created_at, last_seen, display_name)
VALUES ($1, $2, $2, $3)
RETURNING guest_id, display_name, created_at, last_seen`, id, now, name).
			Scan(&v.ID, &v.DisplayName, &v.CreatedAt, &v.LastSeen)
		if err != nil {
			return nil, err
		}

		if err := tx.Commit(); err != nil {
			return nil, err
		}

		return &v, nil
	default:
		return nil, err
	}
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
	return count, nil
}
