package db

import (
	"log"
	"time"
)

type Pet struct {
	PetID       string    `sql:"pet_id"`
	Name        string    `sql:"pet_name"`
	createdAt   time.Time `sql:"created_at"`
	updatedAt   time.Time `sql:"updated_at"`
	Visibility  bool      `sql:"visibility"`
	userID      *string   `sql:"user_id"`
	ActiveImage *string   `sql:"active_image"`
}

func (c Client) GetAllPets() ([]Pet, error) {
	var pets []Pet
	query := `SELECT * FROM pet;`

	rows, err := c.db.Query(query)
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

func (c Client) GetPetByID(id string) (*Pet, error) {
	var pet Pet
	query := `SELECT * FROM pet WHERE pet_id = $1;`
	row := c.db.QueryRow(query, id)

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

func (c Client) GetPetImage(id string) string {
	var imageURL string
	query := `SELECT image_url FROM PetImage WHERE image_id = $1;`

	row := c.db.QueryRow(query, id)

	err := row.Scan(&imageURL)
	if err != nil {
		log.Printf("Failed to fetch image: %v\n", err)

		// temporary
		placeholder := "/assets/placeholder.jpg"
		return placeholder
	}

	return imageURL
}
