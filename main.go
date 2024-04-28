package main

import (
	"log"

	"github.com/animeshdas2000/privc/cache"
	"github.com/animeshdas2000/privc/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	redisClient := cache.InitRedis()
	r := gin.Default()
	r.Use(middleware.CacheMiddleware(redisClient))
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "true"})
	})
	r.POST("/tokenize", Tokenize)
	r.POST("/detokenize", Detokenize)
	r.Run()
}
