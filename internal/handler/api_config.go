package handler

import (
	"log"
	"net/http"
	"pet-everyone/internal/config"
	"pet-everyone/internal/db"
	"pet-everyone/internal/middleware"
	"pet-everyone/internal/registry"
	"strings"
)

type APIConfig struct {
	cfg config.AppConfig
}

func (cfg *APIConfig) GetSessionToken(token string) (*db.SessionToken, error) {
	db := cfg.GetDB()
	return db.GetSessionToken(token)
}

func SetupAPI(dbClient db.Client, filepathRoot string, assetsRoot string, port string) *APIConfig {
	cfg := &APIConfig{config.NewConfig(dbClient, filepathRoot, assetsRoot, port)}
	return cfg
}

func SetupMux(cfg *APIConfig, hubRegistry *registry.HubRegistry) *http.ServeMux {

	authMiddleware := middleware.Auth(cfg)

	log.Printf("filepathRoot: %s", cfg.GetFilepathRoot())
	mux := http.NewServeMux()

	appFs := http.FileServer(http.Dir(cfg.GetFilepathRoot()))
	assetFs := http.FileServer(http.Dir(cfg.GetAssetsRoot()))

	mux.Handle("/app/", http.StripPrefix("/app", blockDirectoryListing(appFs)))

	assetHandler := http.StripPrefix("/assets", blockDirectoryListing(assetFs))
	mux.Handle("/assets/", assetHandler)

	mux.HandleFunc("/", cfg.serveHome)

	mux.HandleFunc("/login", cfg.serveLogin)
	mux.HandleFunc("/signup", cfg.serveSignup)

	mux.HandleFunc("/api/login", cfg.handleLogin)
	mux.HandleFunc("/api/signup", cfg.handleSignup)

	mux.HandleFunc("/pet/{pet_id}", cfg.handlePetConnect)
	mux.HandleFunc("/pet/{pet_id}/ws", func(w http.ResponseWriter, r *http.Request) {
		cfg.servePetWebsocket(w, r, hubRegistry)
	})

	mux.HandleFunc("/pet/create", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, cfg.GetFilepathRoot()+"/pet/create.html")
	})
	mux.Handle("/pet/create/submit", authMiddleware(http.HandlerFunc(cfg.handleCreatePet)))

	return mux
}

func blockDirectoryListing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (c *APIConfig) GetDB() db.Client        { return c.cfg.GetDB() }
func (c *APIConfig) GetPort() string         { return c.cfg.GetPort() }
func (c *APIConfig) GetFilepathRoot() string { return c.cfg.GetFilepathRoot() }
func (c *APIConfig) GetAssetsRoot() string   { return c.cfg.GetAssetsRoot() }
