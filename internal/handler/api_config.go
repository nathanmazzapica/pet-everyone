package handler

import (
	"log"
	"net/http"
	"pet-everyone/internal/db"
	"pet-everyone/internal/registry"
	"strings"
)

type APIConfig struct {
	db           db.Client
	filepathRoot string
	assetsRoot   string
	port         string
}

func SetupAPI(dbClient db.Client, filepathRoot string, assetsRoot string, port string) *APIConfig {
	cfg := &APIConfig{
		dbClient,
		filepathRoot,
		assetsRoot,
		port,
	}
	return cfg
}

func SetupMux(cfg *APIConfig, hubRegistry *registry.HubRegistry) *http.ServeMux {
	log.Printf("filepathRoot: %s", cfg.filepathRoot)
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir(cfg.filepathRoot))

	mux.Handle("/app/", http.StripPrefix("/app", blockDirectoryListing(fs)))

	assetHandler := http.StripPrefix("/assets", http.FileServer(http.Dir(cfg.assetsRoot)))
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
		http.ServeFile(w, r, cfg.filepathRoot+"/pet/create.html")
	})
	mux.HandleFunc("/pet/create/submit", cfg.handleCreatePet)

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
