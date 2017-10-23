package nut

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kapmahc/axe/web"
)

type fmUsersSignIn struct {
}

func postUsersSignIn(l string, c *web.Context) (interface{}, error) {
	return web.H{}, nil
}

func init() {
	Mount(func(rt *mux.Router) {
		unt := rt.PathPrefix("/api/users").Subrouter()
		unt.HandleFunc("/sign-in", JSON(postUsersSignIn)).Methods(http.MethodPost)
	})
}
