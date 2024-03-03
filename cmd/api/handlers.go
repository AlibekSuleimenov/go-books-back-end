package main

import (
	"errors"
	"github.com/alibeksuleimenov/go-books-back-end/internal/data"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"time"
)

// JSONResponse is the type for structuring JSON response
type JSONResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Envelope is a simple wrapper for JSON response
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

	// generate token if user is active
	if user.Active == 0 {
		app.errorJSON(w, errors.New("user is not active"))
		return
	}

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

// Logout logs user out by deleting all tokens from db
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

// AllUsers returns all records from users table
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

// EditUser updates User and saves changes to the users table
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
		u.Active = user.Active

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

// GetUser returns user from db by ID
func (app *Application) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	user, err := app.Models.User.GetByID(userID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, user)
}

// DeleteUser removes user from db
func (app *Application) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		ID int `json:"id"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.Models.User.DeleteByID(requestPayload.ID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := JSONResponse{
		Error:   false,
		Message: "user deleted",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)

}

// LogUserOutAndSetInactive updates user to inactive and deletes all tokens from db
func (app *Application) LogUserOutAndSetInactive(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	user, err := app.Models.User.GetByID(userID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	user.Active = 0

	err = user.Update()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// delete tokens
	err = app.Models.Token.DeleteTokensForUser(userID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := JSONResponse{
		Error:   false,
		Message: "user logged out and set to inactive",
	}

	_ = app.writeJSON(w, http.StatusAccepted, payload)
}

// ValidateToken validates user's token
func (app *Application) ValidateToken(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Token string `json:"token"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	valid := false
	valid, _ = app.Models.Token.ValidToken(requestPayload.Token)

	payload := JSONResponse{
		Error: false,
		Data:  valid,
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}
