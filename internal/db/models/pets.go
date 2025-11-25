package models

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

type Pet struct {
	PetID       string
	Name        string
	createdAt   time.Time
	updatedAt   time.Time
	Visibility  bool
	userID      *string
	ActiveImage *string
}

type PetModel struct {
	DB *sql.DB
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
		err := rows.Scan(&pet.PetID, &pet.Name, &pet.createdAt, &pet.updatedAt, &pet.Visibility, &pet.userID, &pet.ActiveImage)
		if err != nil {
			return nil, err
		}
		pets = append(pets, pet)
	}

	return pets, nil
}

func (p *PetModel) Get(id string) (*Pet, error) {
	var pet Pet
	query := `SELECT * FROM pet WHERE pet_id = $1;`
	row := p.DB.QueryRow(query, id)

	err := row.Scan(
		&pet.PetID,
		&pet.Name,
		&pet.createdAt,
		&pet.updatedAt,
		&pet.Visibility,
		&pet.userID,
		&pet.ActiveImage,
	)

	if err != nil {
		return nil, err
	}

	return &pet, nil
}

func (p *PetModel) GetPetImage(id string) string {
	var imageURL string
	query := `SELECT image_url FROM PetImage WHERE image_id = $1;`

	row := p.DB.QueryRow(query, id)

	err := row.Scan(&imageURL)
	if err != nil {
		log.Printf("Failed to fetch image: %v\n", err)

		// temporary
		placeholder := "/assets/placeholder.jpg"
		return placeholder
	}

	return imageURL
}

func (p *PetModel) CreatePet(pet *Pet) error {
	query := `INSERT INTO pet (pet_id, pet_name, visibility, active_image) VALUES ($1, $2, $3, $4);`

	pet.PetID = uuid.New().String()
	_, err := p.DB.Exec(query, pet.PetID, pet.Name, pet.Visibility, pet.ActiveImage)

	return err
}

func (p *PetModel) CreatePetImage(imageURL string) (string, error) {
	query := `INSERT INTO PetImage (image_url) VALUES ($1) RETURNING image_id;`
	var id string
	err := p.DB.QueryRow(query, imageURL).Scan(&id)

	return id, err
}
