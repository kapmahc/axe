package orm

import "database/sql"

var _db *sql.DB

// Tx transaction
func Tx(f func(*sql.Tx) error) error {
	tx, er := _db.Begin()
	if er != nil {
		return er
	}
	er = f(tx)
	if er == nil {
		return tx.Commit()
	}
	return tx.Rollback()
}

// Open open database
func Open(db *sql.DB) {
	_db = db
}
