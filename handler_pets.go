package main

import (
	"html/template"
	"log"
	"net/http"
	"pet-everyone/internal/websocket"
)

// temporary for testing
type Pet struct {
	ID       int
	Name     string
	ImageURL string
}

func handlePetConnect(w http.ResponseWriter, r *http.Request) {
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
func handlePetWebsocketConnection(w http.ResponseWriter, r *http.Request) {
	petID := r.PathValue("pet_id")
	log.Println("Connecting to pet{", petID, "}")
	test := websocket.NewHub()
	go test.Run()
	websocket.ServeWs(test, w, r)

}
