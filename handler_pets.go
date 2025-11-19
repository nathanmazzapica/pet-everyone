package main

import (
	"log"
	"net/http"
	"pet-everyone/internal/db"
	"pet-everyone/internal/registry"
	"pet-everyone/internal/websocket"
)

func (c *apiConfig) serveHome(w http.ResponseWriter, _ *http.Request) {
	pets, err := c.db.GetAllPets()
	if err != nil {
		respondWithError(w, 500, "unable to load pets", err)
		return
	}
	data := struct {
		Pets []db.Pet
	}{
		Pets: pets,
	}

	if err := render(w, "home", data); err != nil {
		respondWithError(w, 500, "failed to populate template", err)
	}
}

func (c *apiConfig) handlePetConnect(w http.ResponseWriter, r *http.Request) {
	petID := r.PathValue("pet_id")
	log.Println("Handling connection for pet{", petID, "}")

	// Get pet metadata from db
	pet, err := c.db.GetPetByID(petID)
	if err != nil {
		respondWithError(w, 500, "unable to load pet", err)
		return
	}

	image := c.db.GetPetImage(*pet.ActiveImage)

	data := struct {
		PetName  string
		PetID    string
		ImageURL string
	}{
		PetName:  pet.Name,
		PetID:    petID,
		ImageURL: image,
	}

	if err := render(w, "pet", data); err != nil {
		respondWithError(w, 500, "failed to populate template", err)
	}
}

// servePetWebsocket directs the user to the appropriate websocket hub for their chosen pet
func (c *apiConfig) servePetWebsocket(w http.ResponseWriter, r *http.Request, reg *registry.HubRegistry) {
	petID := r.PathValue("pet_id")
	log.Println("Connecting to pet{", petID, "}")

	hub, _ := reg.GetOrCreateHub(petID)
	log.Println(hub)
	websocket.ServeWs(hub, w, r)

}
