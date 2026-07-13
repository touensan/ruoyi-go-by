package dal

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host     string
	Port     int
	Database int
	Password string
}

var Redis *redis.Client

func initRedis(config *RedisConfig) {

	Redis = redis.NewClient(&redis.Options{
		Addr:            config.Host + ":" + strconv.Itoa(config.Port),
		Password:        config.Password,
		DB:              config.Database,
		Protocol:        2,    // 使用 RESP2，兼容 Redis 5/6 以及较旧代理。
		DisableIdentity: true, // 避免向旧 Redis 发送 CLIENT SETINFO。
	})

	_, err := Redis.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
}
