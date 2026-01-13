package handler

import (
	"log"
	"net/http"
	"pet-everyone/cmd/web/application"
	"pet-everyone/internal/websocket"
)

func serveHome(app *application.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		data, err := app.GetAllPets()
		if err != nil {
			app.RespondWithError(w, 500, "unable to load pets", err)
			return
		}

		if err := render(w, "home", data); err != nil {
			app.RespondWithError(w, 500, "failed to populate template", err)
		}

	})
}

func handlePetConnect(app *application.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		petID := r.PathValue("pet_id")
		log.Println("Handling connection for pet{", petID, "}")

		petData, err := app.GetPetData(petID)
		if err != nil {
			app.RespondWithError(w, 500, "unable to load pet data", err)
			return
		}

		if err := render(w, "pet", petData); err != nil {
			app.RespondWithError(w, 500, "failed to populate template", err)
		}
	})
}

func servePetWebsocket(app *application.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		petID := r.PathValue("pet_id")

		// Resolve identity via cookies (session preferred, guest allowed). Validation is handled here for now.
		var userID string

		if sessCookie, err := r.Cookie("session_token"); err == nil {
			sess, serr := app.SessionTokenModel().Get(sessCookie.Value)
			if serr == nil && !sess.IsExpired() {
				userID = sess.UserID
			}
		}

		if userID == "" {
			if guestCookie, err := r.Cookie("guest_id"); err == nil {
				// TODO: validate guest exists before allowing WS access
				userID = guestCookie.Value
			}
		}

		if userID == "" {
			app.RespondWithError(w, http.StatusUnauthorized, "missing authentication", nil)
			return
		}

		hub, _ := app.GetRegistry().GetOrCreateHub(petID)
		websocket.ServeWs(hub, userID, w, r)
		app.Logger().Info("Websocket connection established", "pet_id", petID)
	})
}
