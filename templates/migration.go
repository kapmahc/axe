package migrations

import (
	"github.com/go-pg/migrations"
	log "github.com/sirupsen/logrus"
)

func init() {
	migrations.Register(func(db migrations.DB) error {
		log.Info("migrate {{.Version}}_{{.Name}}")
		_, err := db.Exec(`CREATE TABLE t1()`)
		return err
	}, func(db migrations.DB) error {
		log.Info("rollback {{.Version}}_{{.Name}}")
		_, err := db.Exec(`DROP TABLE t1`)
		return err
	})
}
