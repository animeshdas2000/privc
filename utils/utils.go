package utils

import (
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

type Response struct {
	Success      bool        `json:"success"`
	ErrorMessage string      `json:"error_message,omitempty"`
	Data         interface{} `json:"data,omitempty"`
}

type TokenRequestPayload struct {
	Id   string            `json:"id"`
	Data map[string]string `json:"data"`
}

func GenerateAESKey() ([]byte, error) {
	key := make([]byte, 16)
	_, err := rand.Read(key)
	return key, err
}

func ReadEnvironmentVariables(key string) string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Println(err)
	}
	fmt.Print(dir)
	envs, err := godotenv.Read(filepath.Join(dir, ".env"))
	if err != nil {
		log.Fatal(err)
	}
	return envs[key]
}

func GetRedisClientFromCtx(c *gin.Context) *redis.Client {
	redisClient, exists := c.Get("redis_client")
	if !exists {
		log.Printf("redis not found")
		return nil
	}
	return redisClient.(*redis.Client)
}
