package application

import (
	"encoding/json"
	"log"
	"net/http"
)

func (app *Config) RespondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		app.logger.Error("Responding with 5XX error:", "message", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	app.RespondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func (app *Config) RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		app.logger.Error("Error marshalling JSON:", "error", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
	_, err = w.Write(dat)
	if err != nil {
		app.logger.Error("Error writing JSON:", "error", err)
		w.WriteHeader(500)
		return
	}
}
