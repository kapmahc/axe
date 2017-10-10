package settings

import "database/sql"

var _db *sql.DB

// Open open
func Open(db *sql.DB) {
	_db = db
}
