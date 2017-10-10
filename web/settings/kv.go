package settings

import (
	"bytes"
	"database/sql"
	"encoding/gob"

	"github.com/kapmahc/axe/web/orm"
	"github.com/kapmahc/axe/web/security"
)

// Set set
func Set(key string, obj interface{}, encode bool) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(obj)
	if err != nil {
		return err
	}
	var val []byte
	if encode {
		if val, err = security.Encrypt(buf.Bytes()); err != nil {
			return err
		}
	} else {
		val = buf.Bytes()
	}

	if err := orm.Tx(func(tx *sql.Tx) error {
		var id uint
		err := tx.QueryRow(orm.Q("settings.get-id"), key).Scan(&id)
		if err == sql.ErrNoRows {
			_, err = tx.Exec(orm.Q("settings.insert"), key, val, encode)
		} else if err == nil {
			_, err = tx.Exec(orm.Q("i18n.update"), id, val, encode)
		}
		return err
	}); err != nil {
		return err
	}
	return nil
}

// Get get
func Get(key string, obj string) error {
	var val []byte
	var encode bool
	if err := _db.QueryRow(orm.Q("settings.get-by-key"), key).Scan(&val, &encode); err != nil {
		return err
	}

	var buf bytes.Buffer
	dec := gob.NewDecoder(&buf)

	if encode {
		vl, er := security.Decrypt(val)
		if er != nil {
			return er
		}
		buf.Write(vl)
	} else {
		buf.Write(val)
	}

	return dec.Decode(obj)
}
