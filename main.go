package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
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

	mux := http.NewServeMux()
	appHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app", appHandler)

	assetHandler := http.StripPrefix("/assets", http.FileServer(http.Dir(assetsRoot)))
	mux.Handle("/assets/", assetHandler)

	mux.HandleFunc("/pet/{pet_id}", handlePetConnect)
	mux.HandleFunc("/pet/{pet_id}/ws", handlePetWebsocketConnection)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on http://localhost:%s/app\n", port)
	log.Fatal(srv.ListenAndServe())

}
