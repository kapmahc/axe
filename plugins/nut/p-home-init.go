package nut

import (
	"encoding/base64"

	"github.com/SermoDigital/jose/crypto"
	"github.com/facebookgo/inject"
	"github.com/gorilla/mux"
	"github.com/kapmahc/axe/web"
	"github.com/spf13/viper"
)

// Init init beans
func (p *HomePlugin) Init(g *inject.Graph) error {
	db, err := p.openDB()
	if err != nil {
		return err
	}
	secret, err := base64.StdEncoding.DecodeString(viper.GetString("secret"))
	if err != nil {
		return err
	}

	security, err := web.NewSecurity(secret)
	if err != nil {
		return err
	}

	i18n, err := web.NewI18n("locales", db)
	if err != nil {
		return err
	}
	jobber, err := p.openJobber()
	if err != nil {
		return err
	}
	redis := p.openRedis()

	return g.Provide(
		&inject.Object{Value: db},
		&inject.Object{Value: redis},
		&inject.Object{Value: security},
		&inject.Object{Value: i18n},
		&inject.Object{Value: jobber},
		&inject.Object{Value: mux.NewRouter()},
		&inject.Object{Value: web.NewCache(redis, "cache://")},
		&inject.Object{Value: web.NewSettings(db, security)},
		&inject.Object{Value: web.NewJwt(secret, crypto.SigningMethodHS512)},
		&inject.Object{Value: p.openWrapper(secret)},
	)
}
