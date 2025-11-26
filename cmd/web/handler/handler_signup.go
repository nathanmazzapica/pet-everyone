package handler

import (
	"encoding/json"
	"net/http"
	"pet-everyone/cmd/web/application"
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

		// temp
		type successResponse struct {
			Message string `json:"message"`
		}

		app.RespondWithJSON(w, http.StatusCreated, successResponse{
			Message: "Successfully created user",
		})

	})
}
