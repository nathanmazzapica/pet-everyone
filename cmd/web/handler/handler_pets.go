package handler

import (
	"log"
	"net/http"
	"pet-everyone/cmd/web/application"
	"pet-everyone/internal/websocket"
)

func serveHome(app *application.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		data, err := app.GetAllPets()
		if err != nil {
			app.RespondWithError(w, 500, "unable to load pets", err)
			return
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

		petData, err := app.GetPetData(petID)
		if err != nil {
			app.RespondWithError(w, 500, "unable to load pet data", err)
			return
		}

		if err := render(w, "pet", petData); err != nil {
			app.RespondWithError(w, 500, "failed to populate template", err)
		}
	})
}

func servePetWebsocket(app *application.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		petID := r.PathValue("pet_id")

		hub, _ := app.GetRegistry().GetOrCreateHub(petID)
		websocket.ServeWs(hub, w, r)
		app.Logger().Info("Websocket connection established", "pet_id", petID)
	})
}
