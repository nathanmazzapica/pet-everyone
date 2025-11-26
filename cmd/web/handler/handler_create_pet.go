package handler

import (
	"log"
	"net/http"
	"pet-everyone/cmd/web/application"
)

func handleCreatePet(app *application.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			app.RespondWithError(w, http.StatusMethodNotAllowed, "forbidden method", nil)
			return
		}

		const maxMemory = 10 << 20

		err := r.ParseMultipartForm(maxMemory)
		if err != nil {
			app.RespondWithError(w, http.StatusInternalServerError, "unable to parse multipart form", err)
			return
		}

		// TODO: add validation & moderation
		petName := r.FormValue("petName")

		if petName == "" {
			app.RespondWithError(w, http.StatusBadRequest, "must provide petName", nil)
			return
		}

		log.Println("attempting to create new pet", petName)

		// image
		file, header, err := r.FormFile("petImageFile")
		if err != nil {
			app.RespondWithError(w, http.StatusInternalServerError, "could not read pet image", nil)
		}
		defer file.Close()

		pet, err := app.CreateNewPet(petName, file, header)
		if err != nil {
			log.Println("error creating pet:", err)
			app.RespondWithError(w, http.StatusInternalServerError, "Unable to create pet", err)
			return
		}

		app.RespondWithJSON(w, http.StatusCreated, pet)
	})
}
