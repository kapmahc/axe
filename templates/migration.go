package migrations

import (
	"fmt"

	"github.com/go-pg/migrations"
)

func init() {
	migrations.Register(func(db migrations.DB) error {
		fmt.Println("migrate {{.Version}}_{{.Name}}")
		_, err := db.Exec(`CREATE TABLE t1()`)
		return err
	}, func(db migrations.DB) error {
		fmt.Println("rollback {{.Version}}_{{.Name}}")
		_, err := db.Exec(`DROP TABLE t1`)
		return err
	})
}
