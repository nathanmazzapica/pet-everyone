package handler

import (
	"log"
	"net/http"
	"pet-everyone/cmd/web/application"
	"pet-everyone/internal/db"
	"pet-everyone/internal/websocket"
)

func serveHome(app *application.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		pets, err := app.GetDB().GetAllPets()
		if err != nil {
			application.RespondWithError(w, 500, "unable to load pets", err)
			return
		}
		data := struct {
			Pets []db.Pet
		}{
			Pets: pets,
		}

		if err := render(w, "home", data); err != nil {
			application.RespondWithError(w, 500, "failed to populate template", err)
		}

	})
}

func handlePetConnect(app *application.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		petID := r.PathValue("pet_id")
		log.Println("Handling connection for pet{", petID, "}")

		db := app.GetDB()

		// Get pet metadata from db
		pet, err := db.GetPetByID(petID)
		if err != nil {
			application.RespondWithError(w, 500, "unable to load pet", err)
			return
		}

		image := db.GetPetImage(*pet.ActiveImage)

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
			application.RespondWithError(w, 500, "failed to populate template", err)
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
