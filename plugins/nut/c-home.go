package nut

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kapmahc/axe/web"
)

func getInstall(l string, h web.H, c *web.Context) error {
	return nil
}

type fmInstall struct {
	Title                string
	Subhead              string
	Name                 string
	Email                string
	Password             string
	PasswordConfirmation string
}

func postInstall(l string, fm interface{}, c *web.Context) error {
	return nil
}

func init() {
	Mount(func(rt *mux.Router) {
		rt.HandleFunc("/install", Application("nut/install", getInstall)).Methods(http.MethodGet)
		rt.HandleFunc("/install", Form("/install", "/install", &fmInstall{}, postInstall)).Methods(http.MethodPost)
	})
}
