package cache

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/garyburd/redigo/redis"
)

//Get get from cache
func Get(key string, val interface{}) error {
	c := _redis.Get()
	defer c.Close()
	bys, err := redis.Bytes(c.Do("GET", _prefix+key))
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	dec := gob.NewDecoder(&buf)
	buf.Write(bys)
	return dec.Decode(val)
}

//Set set cache item
func Set(key string, val interface{}, ttl time.Duration) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(val); err != nil {
		return err
	}

	c := _redis.Get()
	defer c.Close()
	_, err := c.Do("SET", _prefix+key, buf.Bytes(), "EX", int(ttl/time.Second))
	return err
}

// Flush clear cache
func Flush() error {
	c := _redis.Get()
	defer c.Close()
	keys, err := redis.Values(c.Do("KEYS", _prefix+"*"))
	if err == nil && len(keys) > 0 {
		_, err = c.Do("DEL", keys...)
	}
	return err
}

// Keys cache keys
func Keys() ([]string, error) {
	c := _redis.Get()
	defer c.Close()
	keys, err := redis.Strings(c.Do("KEYS", _prefix+"*"))
	if err != nil {
		return nil, err
	}
	for i := range keys {
		keys[i] = keys[i][len(_prefix):]
	}
	return keys, nil
}
