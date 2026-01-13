package handler

import (
	"encoding/json"
	"net/http"
	"pet-everyone/cmd/web/application"
	"pet-everyone/internal/auth"
	"time"
)

func serveSignup(app *application.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := render(w, "signup", nil)
		if err != nil {
			app.RespondWithError(w, http.StatusInternalServerError, "Unable to render signup page", err)
			return
		}
	})
}

func handleSignup(app *application.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err := decoder.Decode(&params)
		if err != nil {
			app.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
			return
		}

		err = app.Signup(params.Email, params.Password)
		if err != nil {
			app.RespondWithError(w, http.StatusInternalServerError, "Unable to create user", err)
			return
		}

		// TODO: fetch created user/profile, persist session token with user ID
		token := auth.GenerateSessionToken()

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    token,
			Path:     "/",
			Expires:  time.Now().UTC().Add(30 * 24 * time.Hour),
			HttpOnly: true,
			Secure:   false, // TODO: set true in production
			SameSite: http.SameSiteLaxMode,
		})

		type successResponse struct {
			Message string `json:"message"`
		}

		app.RespondWithJSON(w, http.StatusCreated, successResponse{
			Message: "user created",
		})

	})
}
