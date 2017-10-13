package nut

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kapmahc/axe/web"
)

func getInstall(l string, h web.H, c *web.Context) error {
	h[TITLE] = I18N().T(l, "nut.install.title")
	return nil
}

type fmInstall struct {
	Title                string `form:"title" validate:"required"`
	Subhead              string `form:"subhead" validate:"required"`
	Name                 string `form:"name" validate:"required"`
	Email                string `form:"email" validate:"email"`
	Password             string `form:"password" validate:"required"`
	PasswordConfirmation string `form:"passwordConfirmation" validate:"eqfield=Password"`
}

func postInstall(l string, o interface{}, c *web.Context) error {
	fm := o.(*fmInstall)
	tx, err := DB().Begin()
	if err != nil {
		return err
	}

	i18n := I18N()
	err = i18n.Set(tx, l, "site.title", fm.Title)
	if err == nil {
		err = i18n.Set(tx, l, "site.subhead", fm.Subhead)
	}

	if err == nil {
		return tx.Commit()
	}
	return tx.Rollback()
}

func init() {
	Mount(func(rt *mux.Router) {
		rt.HandleFunc("/install", Application("nut/install", getInstall)).Methods(http.MethodGet)
		rt.HandleFunc("/install", Form("/users/sign-in", "/install", &fmInstall{}, postInstall)).Methods(http.MethodPost)
	})
}
