package middleware

import (
	"context"
	"net/http"
	"pet-everyone/cmd/web/application"
)

// Identity carries resolved request identity.
type Identity struct {
	IsGuest bool
	UserID  string
	GuestID string
}

// identityKey is the context key type for Identity.
type identityKey struct{}

// IdentityFromContext fetches the Identity from context.
func IdentityFromContext(ctx context.Context) (Identity, bool) {
	id, ok := ctx.Value(identityKey{}).(Identity)
	return id, ok
}

func Auth(app *application.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessCookie, err := r.Cookie("session_token")
			if err != nil {
				app.RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
				return
			}

			sess, err := app.SessionTokenModel().Get(sessCookie.Value)
			if err != nil || sess.IsExpired() {
				app.RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
				return
			}

			id := Identity{IsGuest: false, UserID: sess.UserID}
			ctx := context.WithValue(r.Context(), identityKey{}, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
