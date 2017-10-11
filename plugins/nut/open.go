package nut

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/go-pg/pg"
	"github.com/gorilla/sessions"
	"github.com/kapmahc/axe/web"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

var (
	_db           *pg.DB
	_redis        *redis.Pool
	_sessionStore sessions.Store
	_cache        *web.Cache
	_security     *web.Security
	_settings     *web.Settings
	_jobber       *web.Jobber
	_i18n         *web.I18n
)

// DB db handle
func DB() *pg.DB {
	return _db
}

// CACHE cache handle
func CACHE() *web.Cache {
	return _cache
}

// SETTINGS settings handle
func SETTINGS() *web.Settings {
	return _settings
}

// SECURITY security handle
func SECURITY() *web.Security {
	return _security
}

// REDIS redis pool handle
func REDIS() *redis.Pool {
	return _redis
}

// JOBBER jobber handle
func JOBBER() *web.Jobber {
	return _jobber
}

// I18N i18n handle
func I18N() *web.I18n {
	return _i18n
}

// -------------------------

func openDB() (*pg.DB, error) {
	args := viper.GetStringMap("postgresql")
	opt, err := pg.ParseURL(fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		args["user"],
		args["password"],
		args["host"],
		args["port"],
		args["dbname"],
		args["sslmode"],
	))
	if err != nil {
		return nil, err
	}
	db := pg.Connect(opt)
	return db, nil
}

func openJobber() (*web.Jobber, error) {
	args := viper.GetStringMap("rabbitmq")
	return web.NewJobber(fmt.Sprintf(
		"amqp://%s:%s@%s:%d/%s",
		args["user"],
		args["password"],
		args["host"],
		args["port"],
		args["virtual"],
	), args["queue"].(string))
}

func openRedis() {
	_redis = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, e := redis.Dial(
				"tcp",
				fmt.Sprintf(
					"%s:%d",
					viper.GetString("redis.host"),
					viper.GetInt("redis.port"),
				),
			)
			if e != nil {
				return nil, e
			}
			if _, e = c.Do("SELECT", viper.GetInt("redis.db")); e != nil {
				c.Close()
				return nil, e
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

// Open read config file
func Open(f cli.ActionFunc, beans bool) cli.ActionFunc {
	viper.SetEnvPrefix("axe")
	viper.BindEnv("env")

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	return func(c *cli.Context) error {
		err := viper.ReadInConfig()
		if err != nil {
			return err
		}
		log.Infof("read config from %s", viper.ConfigFileUsed())
		if beans {
			// ------------
			_db, err = openDB()
			if err != nil {
				return err
			}
			// -------------
			openRedis()
			// ------------
			secret, err := base64.StdEncoding.DecodeString(viper.GetString("secret"))
			if err != nil {
				return err
			}
			_security, err = web.NewSecurity(secret)
			if err != nil {
				return err
			}
			_sessionStore = sessions.NewCookieStore(secret)
			// ------------
			_cache = web.NewCache(_redis, "cache://")
			_settings = web.NewSettings(_db, _security)
			// ------------
			_i18n, err = web.NewI18n("locales", _db)
			if err != nil {
				return err
			}
			// ------------
			_jobber, err = openJobber()
			if err != nil {
				return err
			}
			// ------------
		}
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
		"theme": "moon",
	})

	viper.SetDefault("secret", base64.StdEncoding.EncodeToString(secret))

	viper.SetDefault("elasticsearch", map[string]interface{}{
		"host": "localhost",
		"port": 9200,
		"ssl":  false,
	})

}
