package cache

import (
	"context"

	"github.com/animeshdas2000/privc/utils"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func InitRedis() *redis.Client {
	dsn := utils.ReadEnvironmentVariables("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: dsn,
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	return rdb
}
