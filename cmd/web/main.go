package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"pet-everyone/cmd/web/application"
	"pet-everyone/cmd/web/handler"
	"pet-everyone/internal/data/cache"
	"pet-everyone/internal/data/db"
	"pet-everyone/internal/data/model"
	"pet-everyone/internal/registry"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	startTime := time.Now()
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	hostURL := os.Getenv("HOST")
	if hostURL == "" {
		log.Println("[WARNING] 'HOST' env variable not set. Defaulting to localhost")
		hostURL = "http://localhost:8082"
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
	defer client.Close()
	log.Println("Connected to db in", time.Since(startTime).Round(time.Millisecond))

	rdbClient := cache.Connect(cache.LoadConfig())
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdbClient.RDB().Ping(ctx).Err(); err != nil {
		log.Fatalf("Error connecting to redis: %v", err)
	}
	defer rdbClient.Close()
	log.Println("Connected to redis in", time.Since(startTime).Round(time.Millisecond))

	app := application.NewConfig(
		hostURL,
		client,
		rdbClient,
		registry.NewHubRegistry(&model.PetModel{DB: client.DB()}),
		filepathRoot,
		assetsRoot,
		port,
	)

	err = app.EnsureAssetsDir()
	if err != nil {
		log.Fatalf("Couldn't create assets directory: %v", err)
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           handler.Routes(app),
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       10 * time.Second,
	}

	log.Println("Started server in", time.Since(startTime).Round(time.Millisecond))
	log.Printf("Serving on %s:%s/app\n", hostURL, port)
	log.Fatal(srv.ListenAndServe())

}
