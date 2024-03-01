package main

import (
	"github.com/alibeksuleimenov/go-books-back-end/internal/data"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
)

// routes generates routes and attaches them to handlers
func (app *Application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Get("/users/login", app.Login)
	mux.Post("/users/login", app.Login)

	mux.Get("/users/all", func(w http.ResponseWriter, r *http.Request) {
		var users data.User
		all, err := users.GetAll()
		if err != nil {
			app.ErrorLog.Println(err)
			return
		}

		app.writeJSON(w, http.StatusOK, all)
	})

	mux.Get("/users/add", func(w http.ResponseWriter, r *http.Request) {
		var user = data.User{
			Email:     "j.doe@mail.com",
			FirstName: "John",
			LastName:  "Doe",
			Password:  "password",
		}

		app.InfoLog.Println("Adding user...")

		id, err := app.Models.User.Insert(user)
		if err != nil {
			app.ErrorLog.Println(err)
			app.errorJSON(w, err, http.StatusForbidden)
			return
		}

		app.InfoLog.Println("Got back id of", id)

		newUser, err := app.Models.User.GetByID(id)
		if err != nil {
			app.ErrorLog.Println(err)
			app.errorJSON(w, err, http.StatusForbidden)
			return
		}

		app.writeJSON(w, http.StatusOK, newUser)
	})

	return mux
}
