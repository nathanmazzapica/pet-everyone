package handler

import (
	"net/http"
	"pet-everyone/cmd/web/application"
	"time"
)

func handleLogout(app *application.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if sessCookie, err := r.Cookie("session_token"); err == nil {
			_ = app.SessionTokenModel().Delete(sessCookie.Value)
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   false, // TODO: set true in production
			SameSite: http.SameSiteLaxMode,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "session_present",
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			HttpOnly: false,
			Secure:   false, // TODO: set true in production
			SameSite: http.SameSiteLaxMode,
		})

		// TODO: optionally issue/refresh guest_id cookie here once guest flow is wired

		app.RespondWithJSON(w, http.StatusOK, map[string]string{
			"message": "logged out",
		})

	})
}
