package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"pet-everyone/cmd/web/application"
	"pet-everyone/internal/auth"
	"pet-everyone/internal/data/model"
	"pet-everyone/internal/displayname"
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
			Email       string `json:"email"`
			Password    string `json:"password"`
			DisplayName string `json:"display_name"`
		}

		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err := decoder.Decode(&params)
		if err != nil {
			app.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
			return
		}

		if err := displayname.Validate(params.DisplayName); err != nil {
			app.RespondWithError(w, http.StatusBadRequest, "Invalid display name", err)
			return
		}

		user, err := app.Signup(params.Email, params.Password, params.DisplayName)
		if err != nil {
			status := http.StatusInternalServerError
			if errors.Is(err, model.ErrDisplayNameTaken) {
				status = http.StatusConflict
			}
			app.RespondWithError(w, status, "Unable to create user", err)
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

		http.SetCookie(w, &http.Cookie{
			Name:     "session_present",
			Value:    "1",
			Path:     "/",
			Expires:  time.Now().UTC().Add(30 * 24 * time.Hour),
			HttpOnly: false,
			Secure:   false, // TODO: set true in production
			SameSite: http.SameSiteLaxMode,
		})

		type successResponse struct {
			User    model.User `json:"user"`
			Message string     `json:"message"`
		}

		app.RespondWithJSON(w, http.StatusCreated, successResponse{
			User:    *user,
			Message: "user created",
		})

	})
}
