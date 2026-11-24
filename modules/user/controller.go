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
		res.ResponseClient(c, http.StatusBadRequest, "Login Form not valid", err.Error(), 0	)
		return
	}
	
	user, err, statusError := controller.Service.HandleLoginForm(loginRequest, c)
	if err != nil {
		log.Printf("Login got error: %s\n", err)
		res.ResponseClient(c, statusError, "Login failed", err.Error(), 0)
		return
	}
	res.ResponseClient(c, 200, "Login successful", user, 0)

}

func (controller *Controller) GetListUser(c *gin.Context) {
	service := controller.Service
	users, total, err, statusError := service.GetListUser(c)
	if err != nil {
		log.Printf("Get User got error: %s\n", err)
		res.ResponseClient(c, statusError, "Get User failed", err.Error(), 0)
		return
	}
	res.ResponseClient(c, 200, "Get User successful", users, total)
}

func (controller *Controller) CreateOrUpdateUser(c *gin.Context) {	
	var request CreateOrUpdateUserRequest
	err := c.ShouldBind(&request)
	if err != nil {
		res.ResponseClient(c, http.StatusBadRequest, "Request params not valid", err.Error(), 0)
		return
	}
	service := controller.Service

	if request.Kind == "create" {
		user, err, statusError := service.CreateNewUser(&request, c)
		if err != nil {
			log.Printf("CreateOrUpdateUser got error: %s\n", err)
			res.ResponseClient(c, statusError, "CreateOrUpdateUser failed", err.Error(), 0)
			return
		}
		res.ResponseClient(c, 200, "Create New User successful", user, 0)
	} else if request.Kind == "update" {
		user, err, statusError := service.UpdateUser(&request, c)
		if err != nil {
			log.Printf("CreateOrUpdateUser got error: %s\n", err)
			res.ResponseClient(c, statusError, "CreateOrUpdateUser failed", err.Error(), 0	)
			return
		}
		res.ResponseClient(c, 200, "Update User successful", user, 0)
	}
	
}
