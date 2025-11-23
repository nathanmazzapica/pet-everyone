package handler

import (
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"pet-everyone/cmd/web/application"
	"pet-everyone/internal/db"
)

func handleCreatePet(app *application.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			application.RespondWithError(w, http.StatusMethodNotAllowed, "forbidden method", nil)
			return
		}

		const maxMemory = 10 << 20

		err := r.ParseMultipartForm(maxMemory)
		if err != nil {
			application.RespondWithError(w, http.StatusInternalServerError, "unable to parse multipart form", err)
			return
		}

		// TODO: add validation & moderation
		petName := r.FormValue("petName")

		if petName == "" {
			application.RespondWithError(w, http.StatusBadRequest, "must provide petName", nil)
			return
		}

		log.Println("attempting to create new pet", petName)

		// image
		file, header, err := r.FormFile("petImageFile")
		if err != nil {
			application.RespondWithError(w, http.StatusInternalServerError, "could not read pet image", nil)
		}
		defer file.Close()

		mediaType, _, err := mime.ParseMediaType(header.Header.Get("Content-Type"))
		if err != nil {
			application.RespondWithError(w, http.StatusInternalServerError, "could not read pet image media type", nil)
		}

		if mediaType != "image/jpeg" && mediaType != "image/png" {
			application.RespondWithError(w, http.StatusBadRequest, "Invalid file type", nil)
			return
		}

		assetPath := application.GetAssetPath(mediaType)
		assetDiskPath, err := app.GetAssetDiskPath(assetPath)

		if err != nil {
			application.RespondWithError(w, http.StatusInternalServerError, "Unable to get asset disk path", err)
			return
		}

		// #nosec G304 -- assetDiskPath is server-generated and sanitized in getAssetDiskPath
		dst, err := os.Create(assetDiskPath)
		if err != nil {
			application.RespondWithError(w, http.StatusInternalServerError, "Unable to create file", err)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			application.RespondWithError(w, http.StatusInternalServerError, "Unable to save file", err)
			return
		}

		url := app.GetAssetURL(assetPath)

		log.Println(url)

		database := app.GetDB()

		imageID, err := database.CreatePetImage(url)
		log.Println(imageID)
		if err != nil {
			application.RespondWithError(w, http.StatusInternalServerError, "Unable to create pet image", err)
			return
		}

		// TODO: Add user ID, right now for testing we're just nulling it out
		pet := &db.Pet{
			Name:        petName,
			ActiveImage: &imageID,
			Visibility:  true,
		}

		err = database.CreatePet(pet)
		if err != nil {
			application.RespondWithError(w, http.StatusInternalServerError, "Unable to create pet", err)
			return
		}

		application.RespondWithJSON(w, http.StatusCreated, pet)
	})
}
