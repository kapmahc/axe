package nut

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kapmahc/axe/web"
)

func getUsersSignIn(l string, h web.H, c *web.Context) error {
	return nil
}

type fmUsersSignIn struct {
}

func postUsersSignIn(l string, fm interface{}, c *web.Context) error {
	return nil
}

func init() {
	Mount(func(rt *mux.Router) {
		unt := rt.PathPrefix("/users").Subrouter()
		unt.HandleFunc("/sign-in", Application("nut/users/sign-in", getUsersSignIn)).Methods(http.MethodGet)
		unt.HandleFunc("/sign-in", Form("/users/sign-in", "/users/sign-in", &fmUsersSignIn{}, postUsersSignIn)).Methods(http.MethodPost)
	})
}
