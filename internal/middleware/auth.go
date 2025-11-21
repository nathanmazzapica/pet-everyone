package middleware

import (
	"log"
	"net/http"
	"pet-everyone/internal/auth"
	"pet-everyone/internal/config"
	"pet-everyone/internal/utils"
)

func Auth(cfg config.AppConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := auth.GetTokenFromHeader(r.Header)
			if err != nil {
				utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
				return
			}
			log.Println("token:", token)

			if token == "" {
				utils.RespondWithError(w, http.StatusUnauthorized, "You need to be logged in to do that!", nil)
			}

			db := cfg.GetDB()

			dbToken, err := db.GetSessionToken(token)
			if err != nil {
				utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
				return
			}

			// Need to handle clientside too
			if dbToken.IsExpired() {
				utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized", nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
