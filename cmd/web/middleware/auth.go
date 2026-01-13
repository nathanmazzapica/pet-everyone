package middleware

import (
	"net/http"
	"pet-everyone/cmd/web/application"
)

func Auth(app *application.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: resolve identity from cookies (session_token preferred, guest_id allowed for guest-only paths)
			_, err := r.Cookie("session_token")
			if err != nil {
				app.RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
				return
			}

			// TODO: validate session token against SessionTokenModel and attach identity to context

			next.ServeHTTP(w, r)
		})
	}
}
