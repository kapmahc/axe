package nut

import (
	"net/http"
	"time"

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
	Title                string `json:"title" validate:"required"`
	Subhead              string `json:"subhead" validate:"required"`
	Name                 string `json:"name" validate:"required"`
	Email                string `json:"email" validate:"email"`
	Password             string `json:"password" validate:"required"`
	PasswordConfirmation string `json:"passwordConfirmation" validate:"eqfield=Password"`
}

func postInstall(l string, c *web.Context) (interface{}, error) {
	var fm fmInstall
	if err := c.Bind(&fm); err != nil {
		return nil, err
	}

	now := time.Now()
	ip := c.ClientIP()
	i18n := I18N()
	if err := Tx(func(tx *pg.Tx) error {
		cnt, err := tx.Model(&User{}).Count()
		if err != nil {
			return err
		}
		if cnt > 0 {
			return i18n.E(l, "errors.forbidden")
		}
		if err = i18n.Set(tx, l, "site.title", fm.Title); err != nil {
			return err
		}
		if err = i18n.Set(tx, l, "site.subhead", fm.Subhead); err != nil {
			return err
		}
		user, err := AddEmailUser(tx, fm.Name, fm.Email, fm.Password)
		if err != nil {
			return err
		}
		if err = AddLog(tx, user.ID, ip, i18n.T(l, "nut.logs.sign-up")); err != nil {
			return err
		}
		user.ConfirmedAt = &now
		user.UpdatedAt = now
		if _, err = tx.Model(user).Column("confirmed_at", "updated_at").Update(); err != nil {
			return err
		}
		for _, rn := range []string{RoleRoot, RoleAdmin} {
			if err := Allow(tx, user.ID, rn, DefaultResourceType, DefaultResourceID, 50, 0, 0); err != nil {
				return err
			}
			if err := AddLog(tx, user.ID, ip, i18n.T(l, "nut.logs.apply-role", rn, DefaultResourceType, DefaultResourceID)); err != nil {
				return err
			}
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
