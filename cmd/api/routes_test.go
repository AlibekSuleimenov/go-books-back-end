package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"testing"
)

func Test_Routes_Exist(t *testing.T) {
	testRoutes := testApp.routes()
	chiRoutes := testRoutes.(chi.Router)

	routeExists(t, chiRoutes, "/users/login")
	routeExists(t, chiRoutes, "/users/logout")
	routeExists(t, chiRoutes, "/admin/users/get/{id}")
	routeExists(t, chiRoutes, "/admin/users/save")
	routeExists(t, chiRoutes, "/admin/users")
	routeExists(t, chiRoutes, "/admin/users/delete")
}

func routeExists(t *testing.T, routes chi.Router, route string) {
	// assume that the rout doesn't exist
	found := false

	// walk through all registered routes
	_ = chi.Walk(routes, func(method string, foundRoute string, handler http.Handler, middlewares ...func(handler2 http.Handler) http.Handler) error {
		// if route exists, set found to true
		if route == foundRoute {
			found = true
		}

		return nil
	})

	// if route doesn't exist, fire an error
	if !found {
		t.Errorf("did not find %s in routes", route)
	}
}
