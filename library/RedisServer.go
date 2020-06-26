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
		Addr:     config.GMasterConfig.RedisAddr,
		Password: config.GMasterConfig.RedisPassword,
		DB:       config.GMasterConfig.RedisDB,
	})

	GRedisClient = &RedisServer{RedisClient: c}

	if _, err = c.Ping().Result(); err != nil {
		return
	}
	return
}