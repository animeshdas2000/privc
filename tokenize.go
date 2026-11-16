package main

import (
	"log"
	"net/http"

	"github.com/animeshdas2000/privc/config"
	"github.com/animeshdas2000/privc/utils"
	"github.com/gin-gonic/gin"
)

func Tokenize(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := utils.TokenRequestPayload{}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, utils.Response{
				Success:      false,
				ErrorMessage: "invalid request payload",
			})
			return
		}

		tokens := make(map[string]string, len(payload.Data))
		for field, value := range payload.Data {
			token, err := Encrypt(cfg.EncryptionKey, value)
			if err != nil {
				log.Printf("tokenize: encrypt failed: %v", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, utils.Response{
					Success:      false,
					ErrorMessage: "tokenization failed",
				})
				return
			}
			tokens[field] = token
		}

		c.JSON(http.StatusCreated, utils.Response{
			Success: true,
			Data:    tokens,
		})
	}
}
