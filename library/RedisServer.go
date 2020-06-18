package library

import (
	"github.com/go-redis/redis"
)

type RedisServer struct {
	RedisClient *redis.Client
}

var (
	GRedisClient *RedisServer
)

func InitRedisClient() (err error) {
	var (
		c *redis.Client
	)

	c = redis.NewClient(&redis.Options{
		Addr:               "localhost:6379",
		Password:           "",
		DB:                 0,
	})

	GRedisClient = &RedisServer{RedisClient:c}

	if _, err = c.Ping().Result(); err != nil {
		return
	}
	return
}
