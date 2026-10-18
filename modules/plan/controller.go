package plan

import (
	"errors"
	res "maintenance-system-go/controller"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Controller struct {
	Service *Service
}

func (controller *Controller) CreatePlan(c *gin.Context) {
	var request PlanRequest
	err := c.ShouldBind(&request)
	if err != nil { 
		//Clarify which kind of error inside validation process 
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, fieldErr := range validationErrors {
				if fieldErr.Field() == "Customer" &&
					fieldErr.Tag() == "inCustomerList" {
					res.ResponseClient(
						c,
						http.StatusBadRequest,
						"Request params not valid",
						gin.H{
							"field":       "customer",
							"description": "No customer found with this name",
						},
					)
					return
				}

				if fieldErr.Field() == "NumberVehicle" &&
					fieldErr.Tag() == "inNumberVehicleList" {
					res.ResponseClient(
						c,
						http.StatusBadRequest,
						"Request params not valid",
						gin.H{
							"field":       "number_vehicle",
							"description": "No vehicle found with this number",
						},
					)
					return
				}

				if fieldErr.Field() == "DriverCode" &&
					fieldErr.Tag() == "inDriverList" {
					res.ResponseClient(
						c,
						http.StatusBadRequest,
						"Request params not valid",
						gin.H{
							"field":       "driver_code",
							"description": "No driver found with this code",
						},
					)
					return
				}

				
			}
		}
	
		res.ResponseClient(c, http.StatusBadRequest, "Request params not valid", err.Error())
		return
	}
	service := controller.Service
	plan, err, statusCode := service.CreatePlan(&request, c)
	if err != nil {
		res.ResponseClient(c, statusCode, "Failed to create plan", err.Error())
		return
	}
	res.ResponseClient(c, statusCode, "Plan created successfully", plan) 
}