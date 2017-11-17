package nut

import (
	"encoding/base64"

	"github.com/go-pg/pg"
	"github.com/kapmahc/axe/web"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
)

// HomePlugin admin plugin
type HomePlugin struct {
	I18n     *web.I18n     `inject:""`
	Cache    *web.Cache    `inject:""`
	Jwt      *web.Jwt      `inject:""`
	Jobber   *web.Jobber   `inject:""`
	Settings *web.Settings `inject:""`
	Router   *web.Router   `inject:""`
	DB       *pg.DB        `inject:""`
	Dao      *Dao          `inject:""`
	Layout   *Layout       `inject:""`
}

// Mount register
func (p *HomePlugin) Mount() error {
	p.Router.GET("nut.home", "/", p.getHome)
	p.Router.Form("nut.install", "/install", web.APPLICATION, "nut/install", p.getInstall, p.postInstall)
	p.Router.Crud("nut.leave-words", "/leave-words", web.APPLICATION, "leave-words", nil, p.newLeaveWord, p.createLeaveWord, nil, nil, nil, nil)
	return nil
}

func init() {
	web.Register(&HomePlugin{})
	viper.SetDefault("languages", []string{
		language.AmericanEnglish.String(),
		language.SimplifiedChinese.String(),
		language.TraditionalChinese.String(),
	})

	viper.SetDefault("aws", map[string]interface{}{
		"access_key_id":     "change-me",
		"secret_access_key": "change-me",
		"region":            "change-me",
		"bucket":            "change-me",
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

	secret, _ := web.RandomBytes(32)
	viper.SetDefault("server", map[string]interface{}{
		"port":   8080,
		"name":   "www.change-me.com",
		"theme":  "moon",
		"secure": false,
	})

	viper.SetDefault("secret", base64.StdEncoding.EncodeToString(secret))

	viper.SetDefault("elasticsearch", map[string]interface{}{
		"host": "localhost",
		"port": 9200,
	})
}
