package main

import "net/http"

// AuthTokenMiddleware is a middleware to check if user is authenticated
func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := app.models.Token.AuthenticateToken(r)
		if err != nil {
			payload := jsonResponse{
				Error:   true,
				Message: "invalid authentication credentials",
			}

			_ = app.writeJSON(w, http.StatusUnauthorized, payload)
		}

		next.ServeHTTP(w, r)
	})
}
