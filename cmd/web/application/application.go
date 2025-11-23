package application

import (
	"log/slog"
	"os"
	"pet-everyone/internal/db"
	"pet-everyone/internal/registry"
)

type Config struct {
	db           *db.Client
	registry     *registry.HubRegistry
	filepathRoot string
	assetsRoot   string
	port         string
	logger       *slog.Logger
}

func NewConfig(dbClient *db.Client, registry *registry.HubRegistry, filepathRoot string, assetsRoot string, port string) *Config {
	return &Config{
		db:           dbClient,
		registry:     registry,
		filepathRoot: filepathRoot,
		assetsRoot:   assetsRoot,
		port:         port,
		logger:       slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}
}

func (app *Config) Logger() *slog.Logger               { return app.logger }
func (app *Config) GetDB() *db.Client                  { return app.db }
func (app *Config) GetRegistry() *registry.HubRegistry { return app.registry }
func (app *Config) GetPort() string                    { return app.port }
func (app *Config) GetFilepathRoot() string            { return app.filepathRoot }
func (app *Config) GetAssetsRoot() string              { return app.assetsRoot }
