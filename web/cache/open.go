package cache

import (
	"github.com/garyburd/redigo/redis"
)

var (
	_redis  *redis.Pool
	_prefix string
)

// Open open
func Open(pool *redis.Pool, prefix string) {
	_redis = pool
	_prefix = prefix
}
