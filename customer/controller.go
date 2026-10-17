package customer

import (
	res "maintenance-system-go/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)


type Controller struct {
	Service *Service
}

func (controller *Controller) CreateOrUpdateCustomer(c *gin.Context) { 
	var request CustomerRequest
	err := c.ShouldBind(&request)
	if err != nil {
		res.ResponseClient(c, http.StatusBadRequest, "Request params not valid", err.Error())
		return
	}
	service := controller.Service

	customer, err, statusError := service.CreateCustomer(&request, c.Request.Context())
	if err != nil {
		res.ResponseClient(c, statusError, "Failed to create customer", err.Error())
		return
	}
	res.ResponseClient(c, http.StatusOK, "Customer created successfully", customer)

}