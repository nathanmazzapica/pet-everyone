package handler

import (
	"log"
	"net/http"
	"pet-everyone/cmd/web/application"
	"pet-everyone/cmd/web/middleware"
	"strings"
	"time"
)

func wrap(next http.Handler, mw ...func(http.Handler) http.Handler) http.Handler {
	for _, m := range mw {
		next = m(next)
	}
	return next
}

func Routes(app *application.Config) *http.ServeMux {
	authMiddleware := middleware.Auth(app)
	limiterMiddleware := middleware.RateLimit(app, 1*time.Second, 100)
	assetsLimiter := middleware.RateLimit(app, 5*time.Millisecond, 100)
	petCreateLimiter := middleware.RateLimit(app, 10*time.Minute, 1)

	log.Printf("filepathRoot: %s", app.GetFilepathRoot())
	mux := http.NewServeMux()

	appFs := http.FileServer(http.Dir(app.GetFilepathRoot()))
	assetFs := http.FileServer(http.Dir(app.GetAssetsRoot()))

	mux.Handle("GET /app/", http.StripPrefix("/app", blockDirectoryListing(appFs)))

	assetHandler := http.StripPrefix("/assets", blockDirectoryListing(assetFs))
	mux.Handle("GET /assets/", assetsLimiter(assetHandler))

	mux.Handle("GET /", limiterMiddleware(serveHome(app)))

	mux.Handle("GET /login", serveLogin(app))
	mux.Handle("GET /signup", serveSignup(app))

	mux.Handle("GET /pet/{pet_id}", handlePetConnect(app))
	mux.Handle("GET /pet/{pet_id}/ws", servePetWebsocket(app))

	mux.HandleFunc("GET /pet/create", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, app.GetFilepathRoot()+"/pet/create.html")
	})

	//mux.Handle("POST /api/login", handleLogin(app))
	wrapped := wrap(handleLogin(app), authMiddleware, limiterMiddleware)
	mux.Handle("POST /api/login", wrapped)

	mux.Handle("POST /api/signup", handleSignup(app))
	mux.Handle("POST /api/logout", handleLogout(app))
	mux.Handle("POST /api/guest/recover", handleGuestRecover(app))
	mux.Handle("GET /api/display-name", handleDisplayNameGet(app))
	mux.Handle("PATCH /api/display-name", authMiddleware(handleDisplayNameSet(app)))

	wrappedPetCreate := wrap(handleCreatePet(app), limiterMiddleware, authMiddleware, petCreateLimiter)
	mux.Handle("POST /api/create", wrappedPetCreate)
	//mux.Handle("POST /api/create", authMiddleware(handleCreatePet(app)))
	mux.Handle("GET /api/pet/{pet_id}/count", handlePersonalPetCount(app))

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
