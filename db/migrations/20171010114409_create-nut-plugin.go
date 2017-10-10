package migrations

import (
	"github.com/go-pg/migrations"
	log "github.com/sirupsen/logrus"
)

func init() {
	migrations.Register(func(db migrations.DB) error {
		log.Info("migrate 20171010114409_create-nut-plugin")
		_, err := db.Exec(`CREATE TABLE t1()`)
		return err
	}, func(db migrations.DB) error {
		log.Info("rollback 20171010114409_create-nut-plugin")
		_, err := db.Exec(`DROP TABLE t1`)
		return err
	})
}
