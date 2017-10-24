package nut

import (
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"log/syslog"
	"path"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/go-pg/pg"
	"github.com/kapmahc/axe/web"
	log "github.com/sirupsen/logrus"
	logrus_syslog "github.com/sirupsen/logrus/hooks/syslog"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

var (
	_db       *pg.DB
	_redis    *redis.Pool
	_cache    *web.Cache
	_security *web.Security
	_settings *web.Settings
	_jobber   *web.Jobber
	_i18n     *web.I18n
	_jwt      *web.Jwt
)

// Tx database transaction
func Tx(f func(*pg.Tx) error) error {
	tx, err := DB().Begin()
	if err != nil {
		return err
	}
	err = f(tx)
	if err == nil {
		tx.Commit()
	} else {
		log.Error(err)
		tx.Rollback()
	}
	return err
}

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

// JWT jwt handle
func JWT() *web.Jwt {
	return _jwt
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

	db.OnQueryProcessed(func(evt *pg.QueryProcessedEvent) {
		query, err := evt.FormattedQuery()
		if err != nil {
			log.Error(err)
			return
		}
		log.Debugf("%s %s", time.Since(evt.StartTime), query)
	})
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
		if err := viper.ReadInConfig(); err != nil {
			return err
		}
		if viper.GetString("env") == web.PRODUCTION {
			// ----------
			log.SetLevel(log.InfoLevel)
			wrt, err := syslog.New(syslog.LOG_INFO, viper.GetString("server.name"))
			if err != nil {
				return err
			}
			log.AddHook(&logrus_syslog.SyslogHook{Writer: wrt})
		} else {
			log.SetLevel(log.DebugLevel)
		}

		log.Infof("read config from %s", viper.ConfigFileUsed())
		if beans {
			var err error
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
			_jwt = web.NewJwt(secret)
			// ------------
			web.SetContext(
				secret,
				path.Join("themes", viper.GetString("server.theme"), "views"),
				template.FuncMap{
					"fmt": fmt.Sprintf,
					"dtf": func(t time.Time) string {
						return t.Format(time.RFC822)
					},
					"eq": func(a interface{}, b interface{}) bool {
						return a == b
					},
					"dict": func(values ...interface{}) (map[string]interface{}, error) {
						if len(values)%2 != 0 {
							return nil, errors.New("invalid dict call")
						}
						dict := make(map[string]interface{}, len(values)/2)
						for i := 0; i < len(values); i += 2 {
							key, ok := values[i].(string)
							if !ok {
								return nil, errors.New("dict keys must be strings")
							}
							dict[key] = values[i+1]
						}
						return dict, nil
					},
					"t": func(lang, code string, args ...interface{}) string {
						return I18N().T(lang, code, args...)
					},
					"assets_css": func(u string) template.HTML {
						return template.HTML(fmt.Sprintf(`<link type="text/css" rel="stylesheet" href="%s">`, u))
					},
					"assets_js": func(u string) template.HTML {
						return template.HTML(fmt.Sprintf(`<script src="%s"></script>`, u))
					},
				},
				viper.GetString("env") != "production",
			)
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
