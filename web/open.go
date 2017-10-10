package web

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

// OpenRedis open redis pool
func OpenRedis(host string, port, db int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, e := redis.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
			if e != nil {
				return nil, e
			}
			if _, e = c.Do("SELECT", db); e != nil {
				c.Close()
				return nil, e
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

// OpenDB open database
func OpenDB(driver, source string) (*sql.DB, error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
