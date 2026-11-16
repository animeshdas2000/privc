package main

import (
	"log"

	"github.com/animeshdas2000/privc/config"
	"github.com/animeshdas2000/privc/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	if err := newRouter(cfg).Run(); err != nil {
		log.Fatal(err)
	}
}

func newRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// Health check stays public so the platform can probe it.
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Server Running"})
	})

	api := r.Group("/", middleware.RequireAPIKey(cfg.APIKeys))
	api.POST("/tokenize", Tokenize(cfg))
	api.POST("/detokenize", Detokenize(cfg))

	return r
}
