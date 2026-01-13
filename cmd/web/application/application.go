package application

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"pet-everyone/internal/data/cache"
	"pet-everyone/internal/data/db"
	"pet-everyone/internal/data/model"
	"pet-everyone/internal/dto"
	"pet-everyone/internal/registry"
)

type Config struct {
	petModel          *model.PetModel
	userModel         *model.UserModel
	sessionTokenModel *model.SessionTokenModel
	visitorModel      *model.VisitorModel
	registry          *registry.HubRegistry
	filepathRoot      string
	assetsRoot        string
	port              string
	logger            *slog.Logger
}

func NewConfig(pool *db.Client, rdbPool *cache.Client, registry *registry.HubRegistry, filepathRoot string, assetsRoot string, port string) *Config {
	db := pool.DB()
	rdb := rdbPool.RDB()
	cfg := &Config{
		petModel:          &model.PetModel{DB: db, Cache: rdb},
		userModel:         &model.UserModel{DB: db},
		sessionTokenModel: &model.SessionTokenModel{DB: db},
		visitorModel:      &model.VisitorModel{DB: db},
		registry:          registry,
		filepathRoot:      filepathRoot,
		assetsRoot:        assetsRoot,
		port:              port,
		logger:            slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})),
	}

	if cfg.registry != nil {
		cfg.registry.SetResolverDeps(cfg.userModel, cfg.visitorModel)
	}

	return cfg
}

func (app *Config) GetAllPets() (*dto.PetList, error) {
	pets, err := app.petModel.GetAll()
	if err != nil {
		return nil, err
	}

	petDTOs := make([]dto.Pet, len(pets))
	for i, pet := range pets {
		dto := dto.Pet{PetID: pet.PetID, PetName: pet.Name}
		dto.PetCount, err = app.petModel.GetPetCount(&pet.PetID)
		if err != nil {
			dto.PetCount = 0
		}
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

func (app *Config) Signup(email, password, displayName string) (*model.User, error) {
	user, err := app.UserModel().Create(email, password, displayName)
	if err != nil {
		app.logger.Error("failed to create user", "err", err)
		return nil, err
	}

	app.logger.Info("new user created", "email", email, "user_id", user.ID)

	return user, nil
}

func (app *Config) GetUserIDFromToken(token string) (string, error) {
	id, err := app.SessionTokenModel().GetUserID(token)
	if err != nil {
		app.logger.Error("failed to get user id from token", "err", err)
		return "", err
	}
	return id, nil
}

func (app *Config) Logger() *slog.Logger                        { return app.logger }
func (app *Config) PetModel() *model.PetModel                   { return app.petModel }
func (app *Config) UserModel() *model.UserModel                 { return app.userModel }
func (app *Config) SessionTokenModel() *model.SessionTokenModel { return app.sessionTokenModel }
func (app *Config) VisitorModel() *model.VisitorModel           { return app.visitorModel }
func (app *Config) GetRegistry() *registry.HubRegistry          { return app.registry }
func (app *Config) GetPort() string                             { return app.port }
func (app *Config) GetFilepathRoot() string                     { return app.filepathRoot }
func (app *Config) GetAssetsRoot() string                       { return app.assetsRoot }
