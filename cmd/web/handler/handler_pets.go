package handler

import (
	"log"
	"net/http"
	"pet-everyone/cmd/web/application"
	"pet-everyone/internal/db/models"
	"pet-everyone/internal/websocket"
)

func serveHome(app *application.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		pets, err := app.PetModel().GetAll()
		if err != nil {
			app.RespondWithError(w, 500, "unable to load pets", err)
			return
		}
		data := struct {
			Pets []models.Pet
		}{
			Pets: pets,
		}

		if err := render(w, "home", data); err != nil {
			app.RespondWithError(w, 500, "failed to populate template", err)
		}

	})
}

func handlePetConnect(app *application.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		petID := r.PathValue("pet_id")
		log.Println("Handling connection for pet{", petID, "}")

		// Get pet metadata from db
		pet, err := app.PetModel().Get(petID)
		if err != nil {
			app.RespondWithError(w, 500, "unable to load pet", err)
			return
		}

		image := app.PetModel().GetPetImage(*pet.ActiveImage)

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
			app.RespondWithError(w, 500, "failed to populate template", err)
		}
	})
}

func servePetWebsocket(app *application.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		petID := r.PathValue("pet_id")
		log.Println("Serving websocket for pet{", petID, "}")

		hub, _ := app.GetRegistry().GetOrCreateHub(petID)
		log.Println(hub)
		websocket.ServeWs(hub, w, r)
	})
}
