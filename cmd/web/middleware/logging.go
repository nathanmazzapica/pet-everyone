package middleware

import (
	"net/http"
	"pet-everyone/cmd/web/application"
)

func Logging(app *application.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			app.Logger().Info("Request received", "method", r.Method, "url", r.URL)
			next.ServeHTTP(w, r)
		})
	}
}
