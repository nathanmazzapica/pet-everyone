package handler

import (
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"pet-everyone/cmd/web/application"
	"pet-everyone/internal/db/models"
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

		mediaType, _, err := mime.ParseMediaType(header.Header.Get("Content-Type"))
		if err != nil {
			app.RespondWithError(w, http.StatusInternalServerError, "could not read pet image media type", nil)
		}

		if mediaType != "image/jpeg" && mediaType != "image/png" {
			app.RespondWithError(w, http.StatusBadRequest, "Invalid file type", nil)
			return
		}

		assetPath := application.GetAssetPath(mediaType)
		assetDiskPath, err := app.GetAssetDiskPath(assetPath)

		if err != nil {
			app.RespondWithError(w, http.StatusInternalServerError, "Unable to get asset disk path", err)
			return
		}

		// #nosec G304 -- assetDiskPath is server-generated and sanitized in getAssetDiskPath
		dst, err := os.Create(assetDiskPath)
		if err != nil {
			app.RespondWithError(w, http.StatusInternalServerError, "Unable to create file", err)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			app.RespondWithError(w, http.StatusInternalServerError, "Unable to save file", err)
			return
		}

		url := app.GetAssetURL(assetPath)

		log.Println(url)

		imageID, err := app.PetModel().CreatePetImage(url)
		log.Println(imageID)
		if err != nil {
			app.RespondWithError(w, http.StatusInternalServerError, "Unable to create pet image", err)
			return
		}

		// TODO: Add user ID, right now for testing we're just nulling it out
		pet := &models.Pet{
			Name:        petName,
			ActiveImage: &imageID,
			Visibility:  true,
		}

		err = app.PetModel().CreatePet(pet)
		if err != nil {
			app.RespondWithError(w, http.StatusInternalServerError, "Unable to create pet", err)
			return
		}

		app.RespondWithJSON(w, http.StatusCreated, pet)
	})
}
