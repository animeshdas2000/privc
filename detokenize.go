package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	"github.com/animeshdas2000/privc/utils"
	"github.com/gin-gonic/gin"
)

func AESDecrypt(key string, iv string, token string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(token)

	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return nil, err
	}

	if len(ciphertext) == 0 || len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("invalid token length")
	}

	mode := cipher.NewCBCDecrypter(block, []byte(iv))
	mode.CryptBlocks(ciphertext, ciphertext)

	return PKCS7UnPadding(ciphertext)
}

// PKCS7UnPadding strips PKCS#7 padding, returning an error instead of
// panicking when the padding is malformed (wrong key or tampered token).
func PKCS7UnPadding(src []byte) ([]byte, error) {
	length := len(src)
	if length == 0 {
		return nil, errInvalidToken
	}
	unpadding := int(src[length-1])
	if unpadding == 0 || unpadding > aes.BlockSize || unpadding > length {
		return nil, errInvalidToken
	}
	for _, b := range src[length-unpadding:] {
		if int(b) != unpadding {
			return nil, errInvalidToken
		}
	}

	return src[:length-unpadding], nil
}

var errInvalidToken = errors.New("invalid token")

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
	encryptionKey := utils.ReadEnvironmentVariables("ENCRYPTION_KEY")
	iv := utils.ReadEnvironmentVariables("IV")

	Fields := TokenRequestPayload.Data
	for i, Field := range Fields {
		//TODO: Write Logic for storing in Persistant storage and comparing
		value, err := AESDecrypt(encryptionKey, iv, Field)
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
