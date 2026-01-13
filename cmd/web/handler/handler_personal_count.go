package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"pet-everyone/cmd/web/application"

	"github.com/google/uuid"
)

func handlePersonalPetCount(app *application.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		petID := r.PathValue("pet_id")
		if _, err := uuid.Parse(petID); err != nil {
			app.RespondWithError(w, http.StatusBadRequest, "invalid pet_id", err)
			return
		}

		// Optional: ensure pet exists
		if _, err := app.PetModel().Get(petID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				app.RespondWithError(w, http.StatusNotFound, "pet not found", err)
				return
			}
			app.RespondWithError(w, http.StatusInternalServerError, "failed to load pet", err)
			return
		}

		personalCount, err := resolvePersonalCount(r, app, petID)
		if err != nil {
			status := http.StatusInternalServerError
			switch {
			case errors.Is(err, errUnauthorized):
				status = http.StatusUnauthorized
			case errors.Is(err, errInvalidIdentity):
				status = http.StatusUnauthorized
			}
			app.RespondWithError(w, status, err.Error(), err)
			return
		}

		app.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
			"pet_id":         petID,
			"personal_count": personalCount,
		})
	})
}

var (
	errUnauthorized    = errors.New("missing or invalid identity")
	errInvalidIdentity = errors.New("invalid identity")
)

func resolvePersonalCount(r *http.Request, app *application.Config, petID string) (int64, error) {
	// Try session user first
	if sessCookie, err := r.Cookie("session_token"); err == nil {
		sess, serr := app.SessionTokenModel().Get(sessCookie.Value)
		if serr == nil && !sess.IsExpired() {
			userUUID, parseErr := uuid.Parse(sess.UserID)
			if parseErr != nil {
				return 0, errInvalidIdentity
			}
			count, qerr := app.UserModel().GetPetCountByPetID(&userUUID, &petID)
			return count, qerr
		}
	}

	// Fallback to guest
	guestCookie, gerr := r.Cookie("guest_id")
	if gerr != nil || guestCookie.Value == "" {
		return 0, errUnauthorized
	}
	guestUUID, parseErr := uuid.Parse(guestCookie.Value)
	if parseErr != nil {
		return 0, errInvalidIdentity
	}

	visitorModel := app.VisitorModel()
	count, qerr := visitorModel.GetPetCountByPetID(&guestUUID, &petID)
	return count, qerr
}
