package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"pet-everyone/cmd/web/application"
)

func serveSignup(app *application.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := render(w, "signup", nil)
		if err != nil {
			application.RespondWithError(w, http.StatusInternalServerError, "Unable to render signup page", err)
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
			application.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
			return
		}

		db := app.GetDB()

		user, err := db.CreateUser(params.Email, params.Password)
		if err != nil {
			application.RespondWithError(w, http.StatusInternalServerError, "Unable to create user", err)
			return
		}

		log.Printf("Created new user: %s\n", user.ID.String())

		application.RespondWithJSON(w, http.StatusCreated, user)

	})
}
