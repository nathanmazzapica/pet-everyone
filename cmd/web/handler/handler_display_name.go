package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"pet-everyone/cmd/web/application"
	"pet-everyone/cmd/web/middleware"
	"pet-everyone/internal/data/model"
	"pet-everyone/internal/displayname"

	"github.com/google/uuid"
)

func handleDisplayNameGet(app *application.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if sessCookie, err := r.Cookie("session_token"); err == nil {
			sess, serr := app.SessionTokenModel().Get(sessCookie.Value)
			if serr == nil && !sess.IsExpired() {
				userID, parseErr := uuid.Parse(sess.UserID)
				if parseErr != nil {
					app.RespondWithError(w, http.StatusBadRequest, "Invalid user id", parseErr)
					return
				}
				user, uerr := app.UserModel().Get(&userID)
				if uerr == nil {
					app.RespondWithJSON(w, http.StatusOK, map[string]string{"display_name": user.DisplayName})
					return
				}
			}
		}

		if guestCookie, err := r.Cookie("guest_id"); err == nil && guestCookie.Value != "" {
			guestID, parseErr := uuid.Parse(guestCookie.Value)
			if parseErr != nil {
				app.RespondWithError(w, http.StatusBadRequest, "Invalid guest id", parseErr)
				return
			}

			visitor, verr := app.VisitorModel().Get(guestID)
			switch {
			case verr == nil:
				app.RespondWithJSON(w, http.StatusOK, map[string]string{"display_name": visitor.DisplayName})
				return
			case errors.Is(verr, model.ErrVisitorNotFound):
				app.RespondWithError(w, http.StatusNotFound, "Display name not found", verr)
				return
			default:
				app.RespondWithError(w, http.StatusInternalServerError, "Unable to fetch display name", verr)
				return
			}
		}

		app.RespondWithError(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	})
}

func handleDisplayNameSet(app *application.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		identity, ok := middleware.IdentityFromContext(r.Context())
		if !ok || identity.IsGuest {
			app.RespondWithError(w, http.StatusUnauthorized, "Unauthorized", nil)
			return
		}

		type request struct {
			DisplayName string `json:"display_name"`
		}

		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			app.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
			return
		}

		if err := displayname.Validate(req.DisplayName); err != nil {
			app.RespondWithError(w, http.StatusBadRequest, "Invalid display name", err)
			return
		}

		userID, err := uuid.Parse(identity.UserID)
		if err != nil {
			app.RespondWithError(w, http.StatusBadRequest, "Invalid user id", err)
			return
		}

		name, err := app.UserModel().SetDisplayName(userID, req.DisplayName)
		if err != nil {
			switch {
			case errors.Is(err, model.ErrDisplayNameTaken):
				app.RespondWithError(w, http.StatusConflict, "Display name already in use", err)
			case errors.Is(err, model.ErrUserNotFound):
				app.RespondWithError(w, http.StatusNotFound, "User not found", err)
			default:
				app.RespondWithError(w, http.StatusInternalServerError, "Unable to set display name", err)
			}
			return
		}

		app.RespondWithJSON(w, http.StatusOK, map[string]string{"display_name": name})
	})
}
