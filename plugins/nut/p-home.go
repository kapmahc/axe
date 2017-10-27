package nut

import (
	"encoding/base64"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/kapmahc/axe/web"
	"github.com/spf13/viper"
)

// HomePlugin admin plugin
type HomePlugin struct {
	I18n     *web.I18n     `inject:""`
	Cache    *web.Cache    `inject:""`
	Router   *gin.Engine   `inject:""`
	Jobber   *web.Jobber   `inject:""`
	Settings *web.Settings `inject:""`
	DB       *pg.DB        `inject:""`
	Dao      *Dao          `inject:""`
	Layout   *Layout       `inject:""`
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

	secret, _ := web.RandomBytes(32)
	viper.SetDefault("server", map[string]interface{}{
		"port":  8080,
		"name":  "www.change-me.com",
		"theme": "moon",
	})

	viper.SetDefault("secret", base64.StdEncoding.EncodeToString(secret))

	viper.SetDefault("elasticsearch", map[string]interface{}{
		"host": "localhost",
		"port": 9200,
	})
}
