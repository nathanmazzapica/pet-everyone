package middleware

import (
	"log"
	"net/http"
	"pet-everyone/cmd/web/application"
	"pet-everyone/internal/auth"
)

func Auth(app *application.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := auth.GetTokenFromHeader(r.Header)
			if err != nil {
				app.RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
				return
			}
			log.Println("token:", token)

			dbToken, err := app.SessionTokenModel().Get(token)
			if err != nil {
				app.RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
				return
			}

			// Need to handle clientside too
			if dbToken.IsExpired() {
				app.RespondWithError(w, http.StatusUnauthorized, "Unauthorized", nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
