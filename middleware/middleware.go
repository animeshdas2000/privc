package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func CacheMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("redis_client", redisClient)
		ctx.Next()
	}
}
