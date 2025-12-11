package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Pet struct {
	PetID       string    `json:"pet_id" db:"pet_id" redis:"pet_id"`
	Name        string    `json:"pet_name" db:"pet_name" redis:"pet_name"`
	CreatedAt   time.Time `json:"created_at" db:"created_at" redis:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at" redis:"updated_at"`
	Visibility  bool      `json:"visibility" db:"visibility" redis:"visibility"`
	UserID      *string   `json:"user_id" db:"user_id" redis:"user_id"`
	ActiveImage *string   `json:"active_image" db:"active_image" redis:"active_image"`
}

type PetModel struct {
	DB    *sql.DB
	Cache *redis.Client
}

func (p *PetModel) GetAll() ([]Pet, error) {
	var pets []Pet
	query := `SELECT * FROM pet;`

	rows, err := p.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var pet Pet
		err := rows.Scan(&pet.PetID, &pet.Name, &pet.CreatedAt, &pet.UpdatedAt, &pet.Visibility, &pet.UserID, &pet.ActiveImage)
		if err != nil {
			return nil, err
		}
		pets = append(pets, pet)
	}

	return pets, nil
}

func (p *PetModel) Get(id string) (*Pet, error) {
	var pet Pet

	pet, hit, err := checkCache(id, p.Cache)

	if hit {
		fmt.Println("[CACHE] :", "hit")
		return &pet, err
	} else {
		fmt.Println("[CACHE] :", "pet_id: ", id, " not found in cache", "error:", err)
		fmt.Println("fetching from db")
	}

	query := `SELECT * FROM pet WHERE pet_id = $1;`
	row := p.DB.QueryRow(query, id)

	err = row.Scan(
		&pet.PetID,
		&pet.Name,
		&pet.CreatedAt,
		&pet.UpdatedAt,
		&pet.Visibility,
		&pet.UserID,
		&pet.ActiveImage,
	)

	if err != nil {
		return nil, err
	}

	err = storeInCache(id, &pet, p.Cache)
	if err != nil {
		fmt.Println("failed to store in cache")
	}

	return &pet, err
}

// pet cache key -> (pet:id):w
func checkCache(id string, client *redis.Client) (Pet, bool, error) {
	var pet Pet
	err := client.HGetAll(context.Background(), "pet:"+id).Scan(&pet)
	hit := pet != (Pet{})

	return pet, hit, err
}

func storeInCache(id string, pet *Pet, client *redis.Client) error {
	return client.HSet(context.Background(), "pet:"+id, pet).Err()
}

func (p *PetModel) GetPetImage(id *string) (string, error) {
	var imageURL string
	query := `SELECT image_url FROM PetImage WHERE image_id = $1;`

	row := p.DB.QueryRow(query, id)

	err := row.Scan(&imageURL)
	if err != nil {
		return "", err
	}

	return imageURL, err
}

func (p *PetModel) GetPetCount(id *string) (int, error) {
	var count int
	query := `SELECT COALESCE(SUM(click_count), 0) FROM UserPetsClickCount WHERE pet_id = $1;`

	row := p.DB.QueryRow(query, id)
	err := row.Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, err
}

func (p *PetModel) CreatePet(pet *Pet) error {
	query := `INSERT INTO pet (pet_id, pet_name, visibility, active_image) VALUES ($1, $2, $3, $4);`

	pet.PetID = uuid.New().String()
	_, err := p.DB.Exec(query, pet.PetID, pet.Name, pet.Visibility, pet.ActiveImage)

	if err != nil {
		return err
	}

	return err
}

func (p *PetModel) CreatePetImage(imageURL string) (string, error) {
	query := `INSERT INTO PetImage (image_url) VALUES ($1) RETURNING image_id;`
	var id string
	err := p.DB.QueryRow(query, imageURL).Scan(&id)

	return id, err
}
