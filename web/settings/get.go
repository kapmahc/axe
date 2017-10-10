package settings

import (
	"bytes"
	"encoding/gob"

	"github.com/kapmahc/axe/web/orm"
	"github.com/kapmahc/axe/web/security"
)

// Get get
func Get(key string, obj string) error {
	var val []byte
	var encode bool
	if err := orm.DB().QueryRow(orm.Q("settings.get-by-key"), key).Scan(&val, &encode); err != nil {
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
