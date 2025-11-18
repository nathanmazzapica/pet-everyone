package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"pet-everyone/internal/db"
	"pet-everyone/internal/registry"
	"time"

	"github.com/joho/godotenv"
)

type apiConfig struct {
	db         db.Client
	assetsRoot string
	port       string
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	filepathRoot := os.Getenv("FILEPATH_ROOT")
	if filepathRoot == "" {
		log.Fatal("FILEPATH_ROOT environment variable not set")
	}

	assetsRoot := os.Getenv("ASSETS_ROOT")
	if assetsRoot == "" {
		log.Fatal("ASSETS_ROOT environment variable not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set")
	}

	client, err := db.Connect(db.LoadConfig())
	if err != nil {
		log.Fatalf("Error connecting to db: %v", err)
	}

	cfg := apiConfig{
		db:         client,
		assetsRoot: assetsRoot,
		port:       port,
	}

	err = cfg.ensureAssetsDir()
	if err != nil {
		log.Fatalf("Couldn't create assets directory: %v", err)
	}

	fmt.Println("Connected to db")

	hubRegistry := registry.NewHubRegistry()

	mux := http.NewServeMux()
	appHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", appHandler)

	assetHandler := http.StripPrefix("/assets", http.FileServer(http.Dir(assetsRoot)))
	mux.Handle("/assets/", assetHandler)

	mux.HandleFunc("/", cfg.serveHome)
	mux.HandleFunc("/pet/{pet_id}", cfg.handlePetConnect)
	mux.HandleFunc("/pet/{pet_id}/ws", func(w http.ResponseWriter, r *http.Request) {
		cfg.servePetWebsocket(w, r, hubRegistry)
	})

	mux.HandleFunc("/pet/create", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepathRoot+"/pet/create.html")
	})
	mux.HandleFunc("/pet/create/submit", cfg.handleCreatePet)

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       10 * time.Second,
	}

	log.Printf("Serving on http://localhost:%s/app\n", port)
	log.Fatal(srv.ListenAndServe())

}
