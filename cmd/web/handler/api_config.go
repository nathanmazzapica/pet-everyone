package handler

import (
	"log"
	"net/http"
	"pet-everyone/cmd/web/application"
	"pet-everyone/cmd/web/middleware"
	"strings"
)

func Routes(app *application.Config) *http.ServeMux {

	authMiddleware := middleware.Auth(app)

	log.Printf("filepathRoot: %s", app.GetFilepathRoot())
	mux := http.NewServeMux()

	appFs := http.FileServer(http.Dir(app.GetFilepathRoot()))
	assetFs := http.FileServer(http.Dir(app.GetAssetsRoot()))

	mux.Handle("GET /app/", http.StripPrefix("/app", blockDirectoryListing(appFs)))

	assetHandler := http.StripPrefix("/assets", blockDirectoryListing(assetFs))
	mux.Handle("GET /assets/", assetHandler)

	mux.Handle("GET /", serveHome(app))

	mux.Handle("GET /login", serveLogin(app))
	mux.Handle("GET /signup", serveSignup(app))

	mux.Handle("GET /pet/{pet_id}", handlePetConnect(app))
	mux.Handle("GET /pet/{pet_id}/ws", servePetWebsocket(app))

	mux.HandleFunc("GET /pet/create", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, app.GetFilepathRoot()+"/pet/create.html")
	})

	mux.Handle("POST /api/login", handleLogin(app))
	mux.Handle("POST /api/signup", handleSignup(app))
	mux.Handle("POST /api/create", authMiddleware(handleCreatePet(app)))

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
