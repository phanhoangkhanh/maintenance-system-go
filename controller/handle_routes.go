package controller

import (
	"maintenance-system-go/authenticate"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)
var validate *validator.Validate

func HandleTest(c *gin.Context) {
	ResponseClient(c, 200, "Test route is working", nil)
}

func Login(c *gin.Context) {
	var loginRequest authenticate.LoginRequest
	err := c.ShouldBindJSON(&loginRequest)
	if err != nil {
		ResponseClient(c, http.StatusBadRequest, "Login Form not valid", err.Error())
		return
	} 
	err = authenticate.HandleLoginForm(loginRequest)
	if err != nil {
		ResponseClient(c, http.StatusInternalServerError, "Login failed", err.Error())
		return
	}
	ResponseClient(c, 200, "Login successful", nil)
	
}