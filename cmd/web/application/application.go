package application

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"pet-everyone/internal/db"
	"pet-everyone/internal/db/models"
	"pet-everyone/internal/dto"
	"pet-everyone/internal/registry"
)

type Config struct {
	petModel          *models.PetModel
	userModel         *models.UserModel
	sessionTokenModel *models.SessionTokenModel
	registry          *registry.HubRegistry
	filepathRoot      string
	assetsRoot        string
	port              string
	logger            *slog.Logger
}

func NewConfig(pool *db.Client, registry *registry.HubRegistry, filepathRoot string, assetsRoot string, port string) *Config {
	db := pool.DB()
	return &Config{
		petModel:          &models.PetModel{DB: db},
		userModel:         &models.UserModel{DB: db},
		sessionTokenModel: &models.SessionTokenModel{DB: db},
		registry:          registry,
		filepathRoot:      filepathRoot,
		assetsRoot:        assetsRoot,
		port:              port,
		logger:            slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})),
	}
}

func (app *Config) GetAllPets() (*dto.PetList, error) {
	pets, err := app.petModel.GetAll()
	if err != nil {
		return nil, err
	}

	petDTOs := make([]dto.Pet, len(pets))
	for i, pet := range pets {
		dto := dto.Pet{PetID: pet.PetID, PetName: pet.Name}
		dto.ImageURL = app.GetPetImage(pet.ActiveImage)
		petDTOs[i] = dto
	}

	petListDTO := dto.PetList{Pets: petDTOs}

	return &petListDTO, nil
}

func (app *Config) GetPetData(petID string) (*dto.Pet, error) {
	pet, err := app.petModel.Get(petID)
	if err != nil {
		return nil, err
	}

	imageURL := app.GetPetImage(pet.ActiveImage)

	petDTO := &dto.Pet{
		PetName:  pet.Name,
		PetID:    pet.PetID,
		ImageURL: imageURL,
	}

	return petDTO, nil
}

func (app *Config) GetPetImage(imageID *string) string {
	imageURL, err := app.petModel.GetPetImage(imageID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.Logger().Warn("Pet missing imageURL", "imageID", imageID)
		} else {
			app.Logger().Error("Unable to get pet imageURL", "imageID", imageID, "err", err)
		}
		imageURL = app.placeholderImage()
	}
	return imageURL
}

func (app *Config) Login() error {
	// handle user login
	return fmt.Errorf("not implemented")
}

func (app *Config) Signup(email, password string) error {
	user, err := app.UserModel().Create(email, password)
	if err != nil {
		app.logger.Error("failed to create user", "err", err)
		return err
	}

	app.logger.Info("new user created", "email", email, "user_id", user.ID)

	return nil
}

func (app *Config) Logger() *slog.Logger                         { return app.logger }
func (app *Config) PetModel() *models.PetModel                   { return app.petModel }
func (app *Config) UserModel() *models.UserModel                 { return app.userModel }
func (app *Config) SessionTokenModel() *models.SessionTokenModel { return app.sessionTokenModel }
func (app *Config) GetRegistry() *registry.HubRegistry           { return app.registry }
func (app *Config) GetPort() string                              { return app.port }
func (app *Config) GetFilepathRoot() string                      { return app.filepathRoot }
func (app *Config) GetAssetsRoot() string                        { return app.assetsRoot }
