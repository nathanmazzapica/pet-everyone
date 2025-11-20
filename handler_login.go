package main

import (
	"encoding/json"
	"net/http"
	"pet-everyone/internal/auth"
	"pet-everyone/internal/db"
	"time"
)

func (cfg *apiConfig) serveLogin(w http.ResponseWriter, r *http.Request) {
	err := render(w, "login", nil)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to render login page", err)
	}
}

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		db.User
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	hash, err := cfg.db.GetUserPasswordHashByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get user password hash", err)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, hash)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to compare password hashes", err)
		return
	}

	if !match {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	user, err := cfg.db.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get user", err)
		return
	}

	token := auth.GenerateSessionToken()
	err = cfg.db.SaveSessionToken(token, user.ID.String(), time.Now().UTC().Add(time.Hour*24*30))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to save session token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{User: *user, Token: token})
}
