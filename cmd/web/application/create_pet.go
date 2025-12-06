package application

import (
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"os/exec"
	"pet-everyone/internal/db/models"
	"strings"
)

var (
	errInvalidMediaType = fmt.Errorf("invalid media type")
)

// CreateNewPet creates a new pet entry in the database with the given name and image. Returns the pet ID or an error.
func (app *Config) CreateNewPet(name string, img multipart.File, header *multipart.FileHeader) (string, error) {
	// Stage 1 : Validate image
	mediaType, _, err := mime.ParseMediaType(header.Header.Get("Content-Type"))
	if err != nil {
		return "", err
	}

	if mediaType != "image/jpeg" && mediaType != "image/png" {
		return "", errInvalidMediaType
	}

	ext := MediaTypeToExtension(mediaType)

	// Stage 2 : Process image
	processedImg, err := app.processImage(img, ext)
	if err != nil {
		return "", err
	}
	defer processedImg.Close()

	// Stage 3 : Save processed image to disk
	// TEMPORARY!!!! WILL BE REPLACED WITH S3 LATER

	assetPath := GetAssetPath(mediaType)
	assetDiskPath, err := app.GetAssetDiskPath(assetPath)

	if err != nil {
		return "", err
	}

	// #nosec G304 -- assetDiskPath is server-generated and sanitized in getAssetDiskPath
	dst, err := os.Create(assetDiskPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, processedImg); err != nil {
		return "", err
	}

	imageURL := app.GetAssetURL(assetPath)

	// Stage 4 : Save pet to database

	pet, err := app.createPetInDatabase(name, imageURL)
	if err != nil {
		return "", err
	}

	return pet.PetID, nil
}

func (app *Config) processImage(img multipart.File, ext string) (*os.File, error) {
	ptrn := fmt.Sprintf("upload-*%s", ext)
	tmpFile, err := os.CreateTemp("", ptrn)
	if err != nil {
		return nil, err
	}

	app.logger.Info("Processing image", "tmpFile", tmpFile.Name())

	if _, err := io.Copy(tmpFile, img); err != nil {
		return nil, err
	}

	pythonExecutable := "./scripts/py/.venv/bin/python"
	removeBGScriptPath := "./scripts/py/bg-removal/main.py"

	cmd := exec.Command(pythonExecutable, removeBGScriptPath, tmpFile.Name())

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to process image: %s", err)
	}

	// Explicit file closure and removal to show we are DONE with it at this stage. Easier in my mind than defer.
	tmpFile.Close()
	os.Remove(tmpFile.Name())

	outputPath := strings.TrimSpace(string(out))

	processedImg, err := os.Open(outputPath)
	if err != nil {
		return nil, err
	}

	return processedImg, nil
}

func (app *Config) createPetInDatabase(name string, imageURL string) (*models.Pet, error) {
	imageID, err := app.petModel.CreatePetImage(imageURL)
	if err != nil {
		return nil, err
	}

	// TODO: Add user ID, right now for testing we're just nulling it out'
	pet := &models.Pet{
		Name:        name,
		ActiveImage: &imageID,
		Visibility:  true,
	}

	err = app.petModel.CreatePet(pet)
	if err != nil {
		return nil, err
	}

	return pet, nil
}
