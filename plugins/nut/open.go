package nut

import (
	"encoding/base64"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

// Open read config file and init beans
func Open(f cli.ActionFunc) cli.ActionFunc {
	viper.SetEnvPrefix("axe")
	viper.BindEnv("env")

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	return func(c *cli.Context) error {
		if err := viper.ReadInConfig(); err != nil {
			return err
		}
		log.Infof("read config from %s", viper.ConfigFileUsed())
		return f(c)
	}
}

func init() {
	viper.SetDefault("env", "development")
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
		"port":     "5672",
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
		"port": 8080,
		"name": "www.change-me.com",
	})

	viper.SetDefault("secret", base64.StdEncoding.EncodeToString(secret))

	viper.SetDefault("elasticsearch", map[string]interface{}{
		"host": "localhost",
		"port": 9200,
	})

}
