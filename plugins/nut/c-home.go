package nut

import (
	"net/http"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"github.com/kapmahc/axe/web"
)

func getAPISiteInfo(l string, c *web.Context) (interface{}, error) {
	i18n := I18N()
	// -----------
	langs, err := i18n.Languages()
	if err != nil {
		return nil, err
	}
	data := web.H{"locale": l, "languages": langs}
	// -----------
	for _, k := range []string{"title", "subhead", "keywords", "description", "copyright"} {
		data[k] = i18n.T(l, "site."+k)
	}
	// -----------
	author := web.H{}
	for _, k := range []string{"name", "email"} {
		author[k] = i18n.T(l, "site.author."+k)
	}
	data["author"] = author
	return data, nil
}

type fmInstall struct {
	Title                string `form:"title" validate:"required"`
	Subhead              string `form:"subhead" validate:"required"`
	Name                 string `form:"name" validate:"required"`
	Email                string `form:"email" validate:"email"`
	Password             string `form:"password" validate:"required"`
	PasswordConfirmation string `form:"passwordConfirmation" validate:"eqfield=Password"`
}

func postInstall(l string, c *web.Context) (interface{}, error) {
	var fm fmInstall
	if err := c.Bind(&fm); err != nil {
		return nil, err
	}

	i18n := I18N()
	if err := Tx(func(tx *pg.Tx) error {
		if err := i18n.Set(tx, l, "site.title", fm.Title); err != nil {
			return err
		}
		if err := i18n.Set(tx, l, "site.subhead", fm.Subhead); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return web.H{}, nil
}

func init() {
	Mount(func(rt *mux.Router) {
		api := rt.PathPrefix("/api").Subrouter()
		api.HandleFunc("/site/info", JSON(getAPISiteInfo)).Methods(http.MethodGet)
		api.HandleFunc("/install", JSON(postInstall)).Methods(http.MethodPost)
	})
}
