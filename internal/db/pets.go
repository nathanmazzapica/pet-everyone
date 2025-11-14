package db

import "time"

type Pet struct {
	PetID       string    `sql:"pet_id"`
	Name        string    `sql:"pet_name"`
	createdAt   time.Time `sql:"created_at"`
	updatedAt   time.Time `sql:"updated_at"`
	Visibility  bool      `sql:"visibility"`
	userID      *string   `sql:"user_id"`
	activeImage *string   `sql:"active_image"`
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
		err := rows.Scan(&pet.PetID, &pet.Name, &pet.createdAt, &pet.updatedAt, &pet.Visibility, &pet.userID, &pet.activeImage)
		if err != nil {
			return nil, err
		}
		pets = append(pets, pet)
	}

	return pets, nil
}
