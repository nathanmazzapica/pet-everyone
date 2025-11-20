package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (c *apiConfig) serveSignup(w http.ResponseWriter, r *http.Request) {
	err := render(w, "signup", nil)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to render signup page", err)
	}
}

func (c *apiConfig) handleSignup(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	user, err := c.db.CreateUser(params.Email, params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create user", err)
		return
	}

	log.Printf("Created new user: %s\n", user.ID.String())

	respondWithJSON(w, http.StatusCreated, user)

}
