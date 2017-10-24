package nut

import (
	"encoding/base64"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"github.com/kapmahc/axe/web"
	"github.com/spf13/viper"
)

// HomePlugin admin plugin
type HomePlugin struct {
	I18n   *web.I18n   `inject:""`
	Cache  *web.Cache  `inject:""`
	DB     *pg.DB      `inject:""`
	Router *mux.Router `inject:""`
	Jobber *web.Jobber `inject:""`
}

// Mount register
func (p *HomePlugin) Mount() error {
	return nil
}

func init() {
	web.Register(&HomePlugin{})

	viper.SetDefault("aws", map[string]interface{}{
		"access_key_id":     "change-me",
		"secret_access_key": "change-me",
		"region":            "change-me",
		"bucket_name":       "change-me",
	})

	viper.SetDefault("redis", map[string]interface{}{
		"host": "localhost",
		"port": 6379,
		"db":   8,
	})

	viper.SetDefault("rabbitmq", map[string]interface{}{
		"user":     "guest",
		"password": "guest",
		"host":     "localhost",
		"port":     5672,
		"virtual":  "axe-dev",
		"queue":    "tasks",
	})

	viper.SetDefault("postgresql", map[string]interface{}{
		"host":     "localhost",
		"port":     5432,
		"user":     "postgres",
		"password": "",
		"dbname":   "axe_dev",
		"sslmode":  "disable",
	})

	secret, _ := RandomBytes(32)
	viper.SetDefault("server", map[string]interface{}{
		"port":  8080,
		"name":  "www.change-me.com",
		"theme": "bootstrap",
	})

	viper.SetDefault("secret", base64.StdEncoding.EncodeToString(secret))

	viper.SetDefault("elasticsearch", map[string]interface{}{
		"host": "localhost",
		"port": 9200,
		"ssl":  false,
	})
}

// func getAPISiteInfo(l string, c *web.Context) (interface{}, error) {
// 	i18n := I18N()
// 	// -----------
// 	langs, err := i18n.Languages()
// 	if err != nil {
// 		return nil, err
// 	}
// 	data := web.H{"locale": l, "languages": langs}
// 	// -----------
// 	for _, k := range []string{"title", "subhead", "keywords", "description", "copyright"} {
// 		data[k] = i18n.T(l, "site."+k)
// 	}
// 	// -----------
// 	author := web.H{}
// 	for _, k := range []string{"name", "email"} {
// 		author[k] = i18n.T(l, "site.author."+k)
// 	}
// 	data["author"] = author
// 	return data, nil
// }
//
// type fmInstall struct {
// 	Title                string `json:"title" validate:"required"`
// 	Subhead              string `json:"subhead" validate:"required"`
// 	Name                 string `json:"name" validate:"required"`
// 	Email                string `json:"email" validate:"email"`
// 	Password             string `json:"password" validate:"required"`
// 	PasswordConfirmation string `json:"passwordConfirmation" validate:"eqfield=Password"`
// }
//
// func postInstall(l string, c *web.Context) (interface{}, error) {
// 	var fm fmInstall
// 	if err := c.Bind(&fm); err != nil {
// 		return nil, err
// 	}
//
// 	now := time.Now()
// 	ip := c.ClientIP()
// 	i18n := I18N()
// 	if err := Tx(func(tx *pg.Tx) error {
// 		cnt, err := tx.Model(&User{}).Count()
// 		if err != nil {
// 			return err
// 		}
// 		if cnt > 0 {
// 			return i18n.E(l, "errors.forbidden")
// 		}
// 		if err = i18n.Set(tx, l, "site.title", fm.Title); err != nil {
// 			return err
// 		}
// 		if err = i18n.Set(tx, l, "site.subhead", fm.Subhead); err != nil {
// 			return err
// 		}
// 		user, err := AddEmailUser(tx, fm.Name, fm.Email, fm.Password)
// 		if err != nil {
// 			return err
// 		}
// 		if err = AddLog(tx, user.ID, ip, i18n.T(l, "nut.logs.sign-up")); err != nil {
// 			return err
// 		}
// 		user.ConfirmedAt = &now
// 		user.UpdatedAt = now
// 		if _, err = tx.Model(user).Column("confirmed_at", "updated_at").Update(); err != nil {
// 			return err
// 		}
// 		for _, rn := range []string{RoleRoot, RoleAdmin} {
// 			if err := Allow(tx, user.ID, rn, DefaultResourceType, DefaultResourceID, 50, 0, 0); err != nil {
// 				return err
// 			}
// 			if err := AddLog(tx, user.ID, ip, i18n.T(l, "nut.logs.apply-role", rn, DefaultResourceType, DefaultResourceID)); err != nil {
// 				return err
// 			}
// 		}
// 		return nil
// 	}); err != nil {
// 		return nil, err
// 	}
// 	return web.H{}, nil
// }
//
// func init() {
// 	Mount(func(rt *mux.Router) {
// 		api := rt.PathPrefix("/api").Subrouter()
// 		api.HandleFunc("/site/info", JSON(getAPISiteInfo)).Methods(http.MethodGet)
// 		api.HandleFunc("/install", JSON(postInstall)).Methods(http.MethodPost)
// 	})
// }
