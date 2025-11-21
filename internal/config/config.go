package config

import (
	"pet-everyone/internal/db"
)

type AppConfig interface {
	GetDB() db.Client
	GetFilepathRoot() string
	GetAssetsRoot() string
	GetPort() string
}

type Config struct {
	db           db.Client
	filepathRoot string
	assetsRoot   string
	port         string
}

func NewConfig(dbClient db.Client, filepathRoot string, assetsRoot string, port string) *Config {
	return &Config{
		db:           dbClient,
		filepathRoot: filepathRoot,
		assetsRoot:   assetsRoot,
		port:         port,
	}
}

func (c *Config) GetDB() db.Client        { return c.db }
func (c *Config) GetPort() string         { return c.port }
func (c *Config) GetFilepathRoot() string { return c.filepathRoot }
func (c *Config) GetAssetsRoot() string   { return c.assetsRoot }
