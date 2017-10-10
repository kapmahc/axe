package nut

import (
	"fmt"
	"sync"

	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	_db    *pg.DB
	dbOnce sync.Once
)

// DB open database
func DB() *pg.DB {
	dbOnce.Do(func() {
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
			log.Error(err)
		}
		_db = pg.Connect(opt)
	})
	return _db
}
