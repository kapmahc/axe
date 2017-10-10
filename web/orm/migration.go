package orm

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	sep  = "_"
	up   = "up.sql"
	down = "down.sql"
)

// Migrate migrate
func Migrate() error {
	return Tx(func(tx *sql.Tx) error {
		items, err := versions()
		if err != nil {
			return err
		}
		for _, s := range []string{
			Q("scheme.create-table"),
			Q("scheme.create-index"),
		} {
			if _, e := tx.Exec(s); e != nil {
				return e
			}
		}

		for _, ver := range items {
			var count int
			if err := _db.QueryRow(Q("scheme.exist"), ver).Scan(&count); err != nil {
				return err
			}
			if count > 0 {
				continue
			}
			us, err := ioutil.ReadFile(filepath.Join("db", "migrations", ver, up))
			if err != nil {
				return err
			}
			if _, err = tx.Exec(string(us)); err != nil {
				return err
			}
			if _, err = tx.Exec(Q("scheme.append"), ver); err != nil {
				return err
			}
		}

		return nil
	})
}

// Rollback rollback
func Rollback() error {
	return Tx(func(tx *sql.Tx) error {
		var ver string
		if err := tx.QueryRow(Q("scheme.latest")).Scan(&ver); err != nil {
			return err
		}
		ds, err := ioutil.ReadFile(filepath.Join("db", "migrations", ver, down))
		if err != nil {
			return err
		}
		log.Debug(string(ds))
		if _, err = tx.Exec(string(ds)); err != nil {
			return err
		}
		_, err = tx.Exec(Q("scheme.remove"), ver)
		return err
	})
}

// Version version
func Version() (string, error) {
	var ver string
	err := _db.QueryRow(Q("scheme.latest")).Scan(&ver)
	return ver, err
}

// Generate generate migration files. [db/migrations]
func Generate() error {
	ver := time.Now().Format("20060102150405")
	root := filepath.Join("db", "migrations", ver+sep)
	if err := os.MkdirAll(root, 0700); err != nil {
		return err
	}
	for _, n := range []string{up, down} {
		fn := filepath.Join(root, n)
		log.Infof("generate file", fn)
		fd, er := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
		if er != nil {
			return er
		}
		defer fd.Close()
	}
	return nil
}
