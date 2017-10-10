package nut

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	_cachePrefix = "cache://"
)

// CacheGet get from cache
func CacheGet(key string, val interface{}) error {
	c := Redis().Get()
	defer c.Close()
	bys, err := redis.Bytes(c.Do("GET", _cachePrefix+key))
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	dec := gob.NewDecoder(&buf)
	buf.Write(bys)
	return dec.Decode(val)
}

// CacheSet set cache item
func CacheSet(key string, val interface{}, ttl time.Duration) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(val); err != nil {
		return err
	}

	c := Redis().Get()
	defer c.Close()
	_, err := c.Do("SET", _cachePrefix+key, buf.Bytes(), "EX", int(ttl/time.Second))
	return err
}

// CacheFlush clear cache
func CacheFlush() error {
	c := Redis().Get()
	defer c.Close()
	keys, err := redis.Values(c.Do("KEYS", _cachePrefix+"*"))
	if err == nil && len(keys) > 0 {
		_, err = c.Do("DEL", keys...)
	}
	return err
}

// CacheKeys cache keys
func CacheKeys() ([]string, error) {
	c := Redis().Get()
	defer c.Close()
	keys, err := redis.Strings(c.Do("KEYS", _cachePrefix+"*"))
	if err != nil {
		return nil, err
	}
	for i := range keys {
		keys[i] = keys[i][len(_cachePrefix):]
	}
	return keys, nil
}
