package user

import (
	"log"
	"net/http"

	res "maintenance-system-go/controller"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	Service *Service
}

func (controller *Controller) Login(c *gin.Context) {
	var loginRequest LoginRequest
	err := c.ShouldBindJSON(&loginRequest)
	if err != nil {
		res.ResponseClient(c, http.StatusBadRequest, "Login Form not valid", err.Error())
		return
	}

	user, err, statusError := controller.Service.HandleLoginForm(loginRequest, c)
	if err != nil {
		log.Printf("Login got error: %s\n", err)
		res.ResponseClient(c, statusError, "Login failed", err.Error())
		return
	}
	res.ResponseClient(c, 200, "Login successful", user)

}
