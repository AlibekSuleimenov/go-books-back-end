package main

import (
	"errors"
	"github.com/alibeksuleimenov/go-books-back-end/internal/data"
	"net/http"
	"time"
)

// JSONResponse
type JSONResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Envelope map[string]interface{}

// Login is the handler used to authenticate users
func (app *Application) Login(w http.ResponseWriter, r *http.Request) {
	type credentials struct {
		Username string `json:"email"`
		Password string `json:"password"`
	}

	var creds credentials
	var payload JSONResponse

	err := app.readJSON(w, r, &creds)
	if err != nil {
		app.ErrorLog.Println(err)
		payload.Error = true
		payload.Message = "invalid json supplied, or json missing entirely"
		_ = app.writeJSON(w, http.StatusBadRequest, payload)
	}

	// authenticate
	// search for a user
	user, err := app.Models.User.GetByEmail(creds.Username)
	if err != nil {
		app.errorJSON(w, errors.New("invalid username"))
		return
	}

	// validate password
	validPassword, err := user.PasswordMatches(creds.Password)
	if err != nil || !validPassword {
		app.errorJSON(w, errors.New("invalid password"))
		return
	}

	// generate token
	token, err := app.Models.Token.GenerateToken(user.ID, 24*time.Hour)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// save token into db
	err = app.Models.Token.Insert(*token, *user)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// send back response
	payload = JSONResponse{
		Error:   false,
		Message: "Logged in successfully!",
		Data:    Envelope{"token": token, "user": user},
	}

	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		app.ErrorLog.Println(err)
	}
}

func (app *Application) Logout(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Token string `json:"token"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, errors.New("invalid JSON"))
		return
	}

	err = app.Models.Token.DeleteByToken(requestPayload.Token)
	if err != nil {
		app.errorJSON(w, errors.New("unable to delete token"))
		return
	}

	payload := JSONResponse{
		Error:   false,
		Message: "Logged Out",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Application) AllUsers(w http.ResponseWriter, r *http.Request) {
	var users data.User
	all, err := users.GetAll()
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	payload := JSONResponse{
		Error:   false,
		Message: "success!",
		Data:    Envelope{"users": all},
	}

	app.writeJSON(w, http.StatusOK, payload)
}

func (app *Application) EditUser(w http.ResponseWriter, r *http.Request) {
	var user data.User

	err := app.readJSON(w, r, &user)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if user.ID == 0 {
		// adding a new user
		if _, err := app.Models.User.Insert(user); err != nil {
			app.errorJSON(w, err)
			return
		}
	} else {
		// editing an existing user
		u, err := app.Models.User.GetByID(user.ID)
		if err != nil {
			app.errorJSON(w, err)
			return
		}

		u.Email = user.Email
		u.FirstName = user.FirstName
		u.LastName = user.LastName

		if err := u.Update(); err != nil {
			app.errorJSON(w, err)
			return
		}

		// if password was changed, update password
		if user.Password != "" {
			err := u.ResetPassword(user.Password)
			if err != nil {
				app.errorJSON(w, err)
				return
			}
		}
	}

	payload := JSONResponse{
		Error:   false,
		Message: "Changes saved!",
	}

	_ = app.writeJSON(w, http.StatusAccepted, payload)
}
