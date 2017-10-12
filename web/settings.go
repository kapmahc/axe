package web

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/go-pg/pg"
)

// NewSettings new settings
func NewSettings(db *pg.DB, sec *Security) *Settings {
	return &Settings{db: db, sec: sec}
}

// Setting setting model
type Setting struct {
	tableName struct{} `sql:"settings"`
	ID        uint
	Key       string
	Value     []byte
	Encode    bool
	Updated   time.Time
	Created   time.Time
}

// Settings settings
type Settings struct {
	db  *pg.DB
	sec *Security
}

// Set set
func (p *Settings) Set(tx *pg.Tx, key string, obj interface{}, encode bool) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(obj)
	if err != nil {
		return err
	}
	var val []byte
	if encode {
		if val, err = p.sec.Encrypt(buf.Bytes()); err != nil {
			return err
		}
	} else {
		val = buf.Bytes()
	}

	var it Setting
	now := time.Now()
	err = tx.Model(&it).Column("id").Where("key = ?", key).Select()
	if err == nil {
		it.Updated = now
		it.Value = val
		it.Encode = encode
		_, err = tx.Model(&it).Column("value", "encode", "updated").Update()
	} else if err == pg.ErrNoRows {
		err = tx.Insert(&Setting{
			Key:     key,
			Value:   val,
			Encode:  encode,
			Updated: now,
			Created: now,
		})
	}

	return err
}

// Get get
func (p *Settings) Get(key string, obj interface{}) error {
	var it Setting
	if err := p.db.Model(&it).
		Column("value", "encode").
		Where("key = ?", key).
		Select(); err != nil {
		return err
	}

	var buf bytes.Buffer
	dec := gob.NewDecoder(&buf)

	if it.Encode {
		vl, err := p.sec.Decrypt(it.Value)
		if err != nil {
			return err
		}
		buf.Write(vl)
	} else {
		buf.Write(it.Value)
	}

	return dec.Decode(obj)
}
