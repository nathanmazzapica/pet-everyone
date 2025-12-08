package application

import (
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"os/exec"
	"pet-everyone/internal/db/models"
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
	assetPath := GenerateAssetPath()
	assetDiskPath, err := app.GetAssetDiskPath(assetPath)
	if err != nil {
		return "", err
	}

	err = app.processImage(img, ext, assetDiskPath)
	if err != nil {
		return "", err
	}

	// Stage 4 : Save pet to database
	imageURL := app.GetAssetURL(assetPath)

	pet, err := app.createPetInDatabase(name, imageURL)
	if err != nil {
		return "", err
	}

	return pet.PetID, nil
}

func (app *Config) processImage(img multipart.File, ext string, outputPath string) error {
	ptrn := fmt.Sprintf("upload-*%s", ext)
	tmpFile, err := os.CreateTemp("", ptrn)
	if err != nil {
		return err
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	app.logger.Info("Processing image", "tmpFile", tmpFile.Name())

	if _, err := io.Copy(tmpFile, img); err != nil {
		return err
	}

	pythonExecutable := "./scripts/py/.venv/bin/python"
	removeBGScriptPath := "./scripts/py/bg-removal/main.py"

	cmd := exec.Command(pythonExecutable, removeBGScriptPath, tmpFile.Name(), outputPath)

	_, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to process image: %s", err)
	}

	return nil
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
