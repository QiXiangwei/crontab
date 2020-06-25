package library

import (
	"crontab/config"
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
		Addr:     config.GConfig.RedisAddr,
		Password: config.GConfig.RedisPassword,
		DB:       config.GConfig.RedisDB,
	})

	GRedisClient = &RedisServer{RedisClient: c}

	if _, err = c.Ping().Result(); err != nil {
		return
	}
	return
}