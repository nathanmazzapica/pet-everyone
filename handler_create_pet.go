package main

import (
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"pet-everyone/internal/db"
)

func (cfg *apiConfig) handlerCreatePet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "forbidden method", nil)
		return
	}

	const maxMemory = 10 << 20

	r.ParseMultipartForm(maxMemory)

	// TODO: add validation & moderation
	petName := r.FormValue("petName")

	if petName == "" {
		respondWithError(w, http.StatusBadRequest, "must provide petName", nil)
		return
	}

	log.Println("attempting to create new pet", petName)

	// image
	file, header, err := r.FormFile("petImageFile")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not read pet image", nil)
	}
	defer file.Close()

	mediaType, _, err := mime.ParseMediaType(header.Header.Get("Content-Type"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not read pet image media type", nil)
	}

	if mediaType != "image/jpeg" && mediaType != "image/png" {
		respondWithError(w, http.StatusBadRequest, "Invalid file type", nil)
		return
	}

	assetPath := getAssetPath(mediaType)
	assetDiskPath := cfg.getAssetDiskPath(assetPath)

	dst, err := os.Create(assetDiskPath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create file", err)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to save file", err)
		return
	}

	url := cfg.getAssetURL(assetPath)

	log.Println(url)
	imageID, err := cfg.db.CreatePetImage(url)
	log.Println(imageID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create pet image", err)
		return
	}

	// TODO: Add user ID, right now for testing we're just nulling it out
	pet := &db.Pet{
		Name:        petName,
		ActiveImage: &imageID,
		Visibility:  true,
	}

	err = cfg.db.CreatePet(pet)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create pet", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, pet)

}
