package main

import (
	"encoding/json"
	"net/http"
)

// JSONResponse
type JSONResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

// Login is the handler used to authenticate users
func (app *Application) Login(w http.ResponseWriter, r *http.Request) {
	type credentials struct {
		Username string `json:"email"`
		Password string `json:"password"`
	}

	var creds credentials
	var payload JSONResponse

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// send back error message
		app.ErrorLog.Println("invalid json")

		payload.Error = true
		payload.Message = "invalid json"

		out, err := json.MarshalIndent(payload, "", "\t")
		if err != nil {
			app.ErrorLog.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(out)
		return
	}

	// TODO authenticate
	app.InfoLog.Println(creds.Username, creds.Password)

	// send back response
	payload.Error = false
	payload.Message = "Signed in"

	out, err := json.MarshalIndent(payload, "", "\t")
	if err != nil {
		app.ErrorLog.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}
