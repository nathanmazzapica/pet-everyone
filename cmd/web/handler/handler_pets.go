package handler

import (
	"log"
	"net/http"
	"time"

	"pet-everyone/cmd/web/application"
	"pet-everyone/internal/websocket"

	"github.com/google/uuid"
)

func serveHome(app *application.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ensureGuestCookie(w, r)

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

		ensureGuestCookie(w, r)

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

		var userID string

		if sessCookie, err := r.Cookie("session_token"); err == nil {
			sess, serr := app.SessionTokenModel().Get(sessCookie.Value)
			if serr == nil && !sess.IsExpired() {
				userID = sess.UserID
			}
		}

		if userID == "" {
			// Allow guests; ensure guest cookie exists before upgrade
			userID = ensureGuestCookie(w, r)
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

func ensureGuestCookie(w http.ResponseWriter, r *http.Request) string {
	if c, err := r.Cookie("guest_id"); err == nil && c.Value != "" {
		return c.Value
	}

	guestID := uuid.New().String()
	http.SetCookie(w, &http.Cookie{
		Name:     "guest_id",
		Value:    guestID,
		Path:     "/",
		Expires:  time.Now().Add(180 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   false, // TODO: set true in production
		SameSite: http.SameSiteLaxMode,
	})

	return guestID
}
