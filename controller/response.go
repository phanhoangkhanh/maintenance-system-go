package controller

import "github.com/gin-gonic/gin"


type ResponseStruct struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
}

func ResponseClient(c *gin.Context, statusCode int,  message string, data any)  {
	c.JSON(statusCode, ResponseStruct{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	})
}