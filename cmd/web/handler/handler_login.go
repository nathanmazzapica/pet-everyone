package handler

import (
	"encoding/json"
	"net/http"
	"pet-everyone/cmd/web/application"
	"pet-everyone/internal/auth"
	"pet-everyone/internal/data/model"
	"time"
)

func serveLogin(app *application.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		err := render(w, "login", nil)
		if err != nil {
			app.RespondWithError(w, http.StatusInternalServerError, "Unable to render login page", err)
			return
		}
	})
}

func handleLogin(app *application.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Password string `json:"password"`
			Email    string `json:"email"`
		}

		type response struct {
			model.User
			Message string `json:"message"`
		}

		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err := decoder.Decode(&params)
		if err != nil {
			app.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
			return
		}

		hash, err := app.UserModel().GetPasswordHashByEmail(params.Email)
		if err != nil {
			app.RespondWithError(w, http.StatusInternalServerError, "Unable to get user password hash", err)
			return
		}

		match, err := auth.CheckPasswordHash(params.Password, hash)
		if err != nil {
			app.RespondWithError(w, http.StatusInternalServerError, "Unable to compare password hashes", err)
			return
		}

		if !match {
			app.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials", nil)
			return
		}

		user, err := app.UserModel().GetByEmail(params.Email)
		if err != nil {
			app.RespondWithError(w, http.StatusInternalServerError, "Unable to get user", err)
			return
		}

		token := auth.GenerateSessionToken()
		err = app.SessionTokenModel().Save(token, user.ID.String(), time.Now().UTC().Add(time.Hour*24*30))
		if err != nil {
			app.RespondWithError(w, http.StatusInternalServerError, "Unable to save session token", err)
			return
		}

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

		app.RespondWithJSON(w, http.StatusOK, response{User: *user, Message: "login successful"})

	})
}
