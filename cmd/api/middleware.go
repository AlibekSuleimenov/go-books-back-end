package main

import "net/http"

// AuthTokenMiddleware
func (app *Application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := app.Models.Token.AuthenticateToken(r)
		if err != nil {
			payload := JSONResponse{
				Error:   true,
				Message: "Invalid authentication credentials",
			}

			_ = app.writeJSON(w, http.StatusUnauthorized, payload)
			return
		}
		next.ServeHTTP(w, r)
	})
}
