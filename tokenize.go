package main

import (
	"fmt"
	"net/http"

	"github.com/animeshdas2000/privc/utils"
	"github.com/gin-gonic/gin"
)

type TokenRequestPayload struct {
	Id   string            `json:"id"`
	Data map[string]string `json:"data"`
}

func Tokenize(c *gin.Context) {
	TokenReqPayload := TokenRequestPayload{}
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
		Field[i] = fmt.Sprintf(val + "Helloiowndjenajn")
	}
	response := utils.Response{
		Success: true,
		Data:    Field,
	}
	c.JSON(http.StatusCreated, response)
}
