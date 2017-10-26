package nut

import (
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"path"
	"time"

	"github.com/SermoDigital/jose/crypto"
	"github.com/facebookgo/inject"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/kapmahc/axe/web"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func (p *HomePlugin) openDB() (*pg.DB, error) {
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

func (p *HomePlugin) openJobber() (*web.Jobber, error) {
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

func (p *HomePlugin) openRouter(secret []byte) *gin.Engine {
	if web.MODE() == web.PRODUCTION {
		gin.SetMode(gin.ReleaseMode)
	}

	rt := gin.Default()
	rt.LoadHTMLGlob(path.Join("themes", viper.GetString("server.theme"), "views") + "/**/*")
	rt.SetFuncMap(template.FuncMap{
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
			return p.I18n.T(lang, code, args...)
		},
		"assets_css": func(u string) template.HTML {
			return template.HTML(fmt.Sprintf(`<link type="text/css" rel="stylesheet" href="%s">`, u))
		},
		"assets_js": func(u string) template.HTML {
			return template.HTML(fmt.Sprintf(`<script src="%s"></script>`, u))
		},
	})

	return rt
}

func (p *HomePlugin) openRedis() *redis.Pool {
	return &redis.Pool{
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
		&inject.Object{Value: gin.Default()},
		&inject.Object{Value: web.NewCache(redis, "cache://")},
		&inject.Object{Value: web.NewSettings(db, security)},
		&inject.Object{Value: web.NewJwt(secret, crypto.SigningMethodHS512)},
		&inject.Object{Value: p.openRouter(secret)},
	)
}
