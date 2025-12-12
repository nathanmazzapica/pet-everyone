package dto

type Pet struct {
	PetName  string
	PetID    string
	PetCount int
	ImageURL string
}

type PetList struct {
	Pets []Pet
}
