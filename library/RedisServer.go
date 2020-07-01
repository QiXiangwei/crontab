package library

import (
	"crontab/config"
	"github.com/go-redis/redis"
	"time"
)

type RedisServer struct {
	RedisClient *redis.Client
}

var (
	GRedisServer *RedisServer
)

func InitRedisServer() (err error) {
	var (
		c *redis.Client
	)

	c = redis.NewClient(&redis.Options{
		Addr:     config.GMasterConfig.RedisAddr,
		Password: config.GMasterConfig.RedisPassword,
		DB:       config.GMasterConfig.RedisDB,
	})

	GRedisServer = &RedisServer{RedisClient: c}

	if _, err = c.Ping().Result(); err != nil {
		return
	}
	return
}

func (redisServer *RedisServer) GetCacheData(key string) (data string, err error) {
	var (
		stringCmd *redis.StringCmd
	)

	stringCmd = redisServer.RedisClient.Get(key)
	if data, err = stringCmd.Result(); err != nil {
		return
	}
	return
}

func (redisServer *RedisServer) SetCacheData(key string, data interface{}, exprTime time.Duration) (result string, err error) {
	var (
		statusCmd *redis.StatusCmd
	)

	statusCmd = redisServer.RedisClient.Set(key, data, exprTime)
	if result, err = statusCmd.Result(); err != nil {
		return
	}
	return
}