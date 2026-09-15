package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)
var validate *validator.Validate

func HandleTest(c *gin.Context) {
	ResponseClient(c, 200, "Test route is working", nil)
}

