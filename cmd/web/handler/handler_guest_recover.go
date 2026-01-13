package handler

import (
	"encoding/json"
	"net/http"
	"pet-everyone/cmd/web/application"
	"time"
)

func handleGuestRecover(app *application.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			RecoveryCode string `json:"recovery_code"`
		}

		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			app.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
			return
		}

		if req.RecoveryCode == "" {
			app.RespondWithError(w, http.StatusBadRequest, "Missing recovery_code", nil)
			return
		}

		// TODO: validate recovery_code against Visitor model and rate limit attempts

		http.SetCookie(w, &http.Cookie{
			Name:     "guest_id",
			Value:    req.RecoveryCode,
			Path:     "/",
			Expires:  time.Now().Add(180 * 24 * time.Hour),
			HttpOnly: true,
			Secure:   false, // TODO: set true in production
			SameSite: http.SameSiteLaxMode,
		})

		app.RespondWithJSON(w, http.StatusOK, map[string]string{
			"message": "guest session restored",
		})
	})
}
