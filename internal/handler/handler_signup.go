package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"pet-everyone/internal/utils"
)

func (c *APIConfig) serveSignup(w http.ResponseWriter, r *http.Request) {
	err := render(w, "signup", nil)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Unable to render signup page", err)
	}
}

func (c *APIConfig) handleSignup(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	user, err := c.db.CreateUser(params.Email, params.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Unable to create user", err)
		return
	}

	log.Printf("Created new user: %s\n", user.ID.String())

	utils.RespondWithJSON(w, http.StatusCreated, user)

}
