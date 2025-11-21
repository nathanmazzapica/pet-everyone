package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"pet-everyone/internal/db"
	"pet-everyone/internal/handler"
	"pet-everyone/internal/registry"
	"time"

	"github.com/joho/godotenv"
)

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

	cfg := handler.SetupAPI(client, filepathRoot, assetsRoot, port)

	err = cfg.EnsureAssetsDir()
	if err != nil {
		log.Fatalf("Couldn't create assets directory: %v", err)
	}

	fmt.Println("Connected to db")

	hubRegistry := registry.NewHubRegistry()

	mux := handler.SetupMux(cfg, hubRegistry)

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       10 * time.Second,
	}

	log.Printf("Serving on http://localhost:%s/app\n", port)
	log.Fatal(srv.ListenAndServe())

}
