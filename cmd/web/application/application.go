package application

import (
	"log/slog"
	"os"
	"pet-everyone/internal/db"
	"pet-everyone/internal/db/models"
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
		logger:            slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}
}

func (app *Config) Logger() *slog.Logger                         { return app.logger }
func (app *Config) PetModel() *models.PetModel                   { return app.petModel }
func (app *Config) UserModel() *models.UserModel                 { return app.userModel }
func (app *Config) SessionTokenModel() *models.SessionTokenModel { return app.sessionTokenModel }
func (app *Config) GetRegistry() *registry.HubRegistry           { return app.registry }
func (app *Config) GetPort() string                              { return app.port }
func (app *Config) GetFilepathRoot() string                      { return app.filepathRoot }
func (app *Config) GetAssetsRoot() string                        { return app.assetsRoot }
