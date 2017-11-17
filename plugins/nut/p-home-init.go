package nut

import (
	"errors"
	"fmt"
	"html/template"
	"time"

	"github.com/SermoDigital/jose/crypto"
	"github.com/facebookgo/inject"
	"github.com/garyburd/redigo/redis"
	"github.com/go-pg/pg"
	"github.com/kapmahc/axe/web"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
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

func (p *HomePlugin) openS3() (*web.S3, error) {
	args := viper.GetStringMapString("aws")
	return web.NewS3(args["access_key_id"], args["secret_access_key"], args["region"], args["bucket"])
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

func (p *HomePlugin) openRouter(secret []byte, db *pg.DB, i18n *web.I18n) (*web.Router, error) {

	helpers := template.FuncMap{
		"fmt": fmt.Sprintf,
		"dtf": func(t time.Time) string {
			return t.Format(time.RFC822)
		},
		"eq": func(a interface{}, b interface{}) bool {
			return a == b
		},
		"str2htm": func(s string) template.HTML {
			return template.HTML(s)
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
		"links": func(lng, loc string) ([]Link, error) {
			var items []Link
			if err := db.Model(&items).Column("id", "label", "href", "loc", "sort_order").
				Where("lang = ? AND loc = ?", lng, loc).
				Order("sort_order ASC").
				Select(); err != nil {
				return nil, err
			}
			return items, nil
		},
		"cards": func(lng, loc string) ([]Card, error) {
			var items []Card
			if err := db.Model(&items).Column("id", "title", "summary", "type", "action", "logo", "href", "loc", "sort_order").
				Where("lang = ? AND loc = ?", lng, loc).
				Order("sort_order ASC").
				Select(); err != nil {
				return nil, err
			}
			return items, nil
		},
		"odd": func(v int) bool {
			return v%2 != 0
		},
		"even": func(v int) bool {
			return v%2 == 0
		},
	}

	theme := viper.GetString("server.theme")

	var langs []language.Tag
	for _, l := range Languages() {
		t, e := language.Parse(l)
		if e != nil {
			return nil, e
		}
		langs = append(langs, t)
	}

	return web.NewRouter(viper.GetBool("server.secure"), secret, theme, helpers), nil
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
	secret, err := web.SECRET()
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

	rt, err := p.openRouter(secret, db, i18n)
	if err != nil {
		return err
	}

	s3, err := p.openS3()
	if err != nil {
		return err
	}

	return g.Provide(
		&inject.Object{Value: db},
		&inject.Object{Value: redis},
		&inject.Object{Value: security},
		&inject.Object{Value: i18n},
		&inject.Object{Value: jobber},
		&inject.Object{Value: s3},
		&inject.Object{Value: web.NewCache(redis, "cache://")},
		&inject.Object{Value: web.NewSettings(db, security)},
		&inject.Object{Value: web.NewJwt(secret, crypto.SigningMethodHS512)},
		&inject.Object{Value: rt},
	)
}
