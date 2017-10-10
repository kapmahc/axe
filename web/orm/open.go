package orm

import "database/sql"

var _db *sql.DB

// DB db
func DB() *sql.DB {
	return _db
}

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
func Open(driver, source string) error {
	db, err := sql.Open(driver, source)
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}
	_db = db
	return nil
}
