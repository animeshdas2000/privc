package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/animeshdas2000/privc/utils"
	"github.com/gin-gonic/gin"
)

func AESDecrypt(token string) ([]byte, error) {
	encryptionKey := utils.ReadEnvironmentVariables("ENCRYPTION_KEY")
	iv := utils.ReadEnvironmentVariables("IV")
	ciphertext, err := base64.StdEncoding.DecodeString(token)

	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher([]byte(encryptionKey))

	if err != nil {
		return nil, err
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("block size cant be zero")
	}

	mode := cipher.NewCBCDecrypter(block, []byte(iv))
	mode.CryptBlocks(ciphertext, ciphertext)
	ciphertext = PKCS5UnPadding(ciphertext)

	return ciphertext, nil
}

// PKCS5UnPadding  pads a certain blob of data with necessary data to be used in AES block cipher
func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])

	return src[:(length - unpadding)]
}

func Detokenize(c *gin.Context) {
	TokenRequestPayload := utils.TokenRequestPayload{}
	err := c.ShouldBindJSON(&TokenRequestPayload)
	if err != nil {
		response := utils.Response{
			Success:      false,
			ErrorMessage: "invalid Request Payload",
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	//redisClient := utils.GetRedisClientFromCtx(c)

	Fields := TokenRequestPayload.Data
	for i, Field := range Fields {
		//TODO: Write Logic for storing in Persistant storage and comparing
		value, err := AESDecrypt(Field)
		if err != nil {
			response := utils.Response{
				Success:      false,
				ErrorMessage: fmt.Sprintf("Detokenization failed %v", err.Error()),
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		Fields[i] = string(value)
	}
	response := utils.Response{
		Success: true,
		Data:    Fields,
	}
	c.JSON(http.StatusOK, response)
}
