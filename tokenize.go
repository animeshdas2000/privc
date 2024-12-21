package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/animeshdas2000/privc/utils"
	"github.com/gin-gonic/gin"
)

func AESEncrypt(key string, iv string, plaintext string) (string, error) {

	var plainTextBlock []byte

	length := len(plaintext)

	if length%16 != 0 || length == 0 {
		extendBlock := 16 - (length % 16)
		plainTextBlock = make([]byte, length+extendBlock)
		copy(plainTextBlock[length:], bytes.Repeat([]byte{uint8(extendBlock)}, extendBlock))
	} else {
		plainTextBlock = make([]byte, length)
	}

	copy(plainTextBlock, []byte(plaintext))
	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return "", err
	}

	cipherText := make([]byte, len(plainTextBlock))
	mode := cipher.NewCBCEncrypter(block, []byte(iv))
	mode.CryptBlocks(cipherText, plainTextBlock)
	str := base64.StdEncoding.EncodeToString(cipherText)
	return str, nil
}

func Tokenize(c *gin.Context) {
	TokenReqPayload := utils.TokenRequestPayload{}
	encryptionKey := utils.ReadEnvironmentVariables("ENCRYPTION_KEY")
	iv := utils.ReadEnvironmentVariables("IV")

	err := c.ShouldBindJSON(&TokenReqPayload)
	if err != nil {
		response := utils.Response{
			Success:      false,
			ErrorMessage: err.Error(),
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	redisClient := utils.GetRedisClientFromCtx(c)

	Field := TokenReqPayload.Data

	// Check if the Request Payload is empty
	if len(Field) == 0 {
		response := utils.Response{
			Success:      false,
			ErrorMessage: "invalid Request Payload",
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// Iterate over the Field map
	for i, val := range Field {
		// Check if the string starts strictly with "Field"
		r, _ := regexp.Compile(`^field`)
		matched := r.MatchString(i)
		if !matched {
			response := utils.Response{
				Success:      false,
				ErrorMessage: fmt.Sprintf("invalid Request payload: '%s' key should start with 'field'", i),
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		// Check if the value is empty
		if val == "" {
			response := utils.Response{
				Success:      false,
				ErrorMessage: fmt.Sprintf("invalid Request payload: '%s' value should not be empty", i),
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		// Check if the value is already present in the cache
		exists, err := redisClient.Exists(c, i).Result()
		if err != nil {
			log.Printf("cache miss for %s: %v", i, err)
			continue
		}

		if exists == 1 {
			response := utils.Response{
				Success:      false,
				ErrorMessage: fmt.Sprintf("invalid Request payload: '%s' key already exists", i),
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		// Encryption
		token, err := AESEncrypt(encryptionKey, iv, val)
		if err != nil {
			response := utils.Response{
				Success:      false,
				ErrorMessage: fmt.Sprintf("tokenization failed %v", err.Error()),
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		Field[i] = token
		err = redisClient.Set(c, i, Field[i], 0).Err()
		if err != nil {
			log.Printf("cache err while set %s: %s", i, err)
		}
	}
	response := utils.Response{
		Success: true,
		Data:    Field,
	}
	c.JSON(http.StatusCreated, response)
}
