package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/animeshdas2000/privc/utils"
	"github.com/gin-gonic/gin"
)

func AESEncrypt(key string, iv string, plaintext string) (string, error) {

	var plainTextBlock []byte

	length := len(plaintext)

	if length%16 != 0 {
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

	Field := TokenReqPayload.Data
	for i, val := range Field {
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
	}
	response := utils.Response{
		Success: true,
		Data:    Field,
	}
	c.JSON(http.StatusCreated, response)
}
