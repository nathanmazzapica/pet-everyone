package main

import (
	"html/template"
	"log"
	"net/http"
	"pet-everyone/internal/db"
	"pet-everyone/internal/registry"
	"pet-everyone/internal/websocket"
)

func (c *apiConfig) handleHome(w http.ResponseWriter, r *http.Request) {
	pets, err := c.db.GetAllPets()
	if err != nil {
		respondWithError(w, 500, "unable to load pets", err)
	}
	data := struct {
		Pets []db.Pet
	}{
		Pets: pets,
	}

	tmpl := template.Must(template.ParseFiles("app/templates/home.html"))

	err = tmpl.Execute(w, data)
	if err != nil {
		respondWithError(w, 500, "unable to load pets", err)
	}
}

func (c *apiConfig) handlePetConnect(w http.ResponseWriter, r *http.Request) {
	petID := r.PathValue("pet_id")
	log.Println("Handling connection for pet{", petID, "}")

	// Get pet metadata from db

	data := struct {
		PetName  string
		PetID    string
		ImageURL string
	}{
		PetName:  "not implemented",
		PetID:    petID,
		ImageURL: "/assets/DAISY.png",
	}

	tmpl := template.Must(template.ParseFiles("app/templates/pet.html"))

	err := tmpl.Execute(w, data)
	if err != nil {
		respondWithError(w, 500, "failed to populate template", err)
	}
}

// handlePetWebsocketConnection directs the user to the appropriate websocket hub for their chosen pet
func (c *apiConfig) handlePetWebsocketConnection(w http.ResponseWriter, r *http.Request, reg *registry.HubRegistry) {
	petID := r.PathValue("pet_id")
	log.Println("Connecting to pet{", petID, "}")

	hub := reg.GetOrCreateHub(petID)
	log.Println(hub)
	websocket.ServeWs(hub, w, r)

}
