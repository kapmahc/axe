package web

import (
	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
)

// Tx database transaction
func Tx(db *pg.DB, f func(*pg.Tx) error) error {
	tx, err := db.Begin()
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
