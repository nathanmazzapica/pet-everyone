package dto

type Pet struct {
	PetName  string
	PetID    string
	ImageURL string
}

type PetList struct {
	Pets []Pet
}
