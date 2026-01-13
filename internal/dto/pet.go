package dto

type Pet struct {
	PetName  string
	PetID    string
	PetCount int64
	ImageURL string
}

type PetList struct {
	Pets []Pet
}
