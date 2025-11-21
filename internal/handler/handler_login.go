package handler

import (
	"encoding/json"
	"net/http"
	"pet-everyone/internal/auth"
	"pet-everyone/internal/db"
	"pet-everyone/internal/utils"
	"time"
)

func (c *APIConfig) serveLogin(w http.ResponseWriter, r *http.Request) {
	err := render(w, "login", nil)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Unable to render login page", err)
	}
}

func (c *APIConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
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
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	db := c.GetDB()

	hash, err := db.GetUserPasswordHashByEmail(params.Email)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Unable to get user password hash", err)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, hash)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Unable to compare password hashes", err)
		return
	}

	if !match {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	user, err := db.GetUserByEmail(params.Email)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Unable to get user", err)
		return
	}

	token := auth.GenerateSessionToken()
	err = db.SaveSessionToken(token, user.ID.String(), time.Now().UTC().Add(time.Hour*24*30))
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Unable to save session token", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, response{User: *user, Token: token})
}
