package i18n

import (
	"database/sql"

	"github.com/kapmahc/axe/web/orm"
)

var (
	_locales = make(map[string]string)
)

const (
	sep = "."
	ext = ".ini"
)

// Set set
func Set(lang, code, message string) error {

	if err := orm.Tx(func(tx *sql.Tx) error {
		var id uint
		err := tx.QueryRow(orm.Q("i18n.get-id"), lang, code).Scan(&id)
		if err == sql.ErrNoRows {
			_, err = tx.Exec(orm.Q("i18n.insert"), lang, code, message)
		} else if err == nil {
			_, err = tx.Exec(orm.Q("i18n.update"), id, message)
		}
		return err
	}); err != nil {
		return err
	}
	_locales[key(lang, code)] = message
	return nil
}

func key(lang, code string) string {
	return lang + sep + code
}
