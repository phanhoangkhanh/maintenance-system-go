package vehicle

import (
	"log"
	"net/http"

	res "maintenance-system-go/controller"

	"github.com/gin-gonic/gin"
)


type Controller struct {
	Service *Service
}

func (controller *Controller) CreateOrUpdateVehicle(c *gin.Context) {
	var request CreateOrUpdateVehicleRequest
	
	err := c.ShouldBind(&request)
	if err != nil {
		res.ResponseClient(c, http.StatusBadRequest, "Request params not valid", err.Error())
		return
	}
	if request.Action == "create" {
		vehicle, err, statusError := controller.Service.CreateNewVehicle(&request, c)
		if err != nil {
			log.Printf("CreateOrUpdateVehicle got error: %s\n", err)
			res.ResponseClient(c, statusError, "CreateOrUpdateVehicle failed", err.Error())
			return
		}
		res.ResponseClient(c, 200, "Create New Vehicle successful", vehicle)
	}

}

func (controller *Controller) AttachVehicle(c *gin.Context) {
	var request AttachVehicleRequest
	err := c.ShouldBind(&request)
	if err != nil {
		res.ResponseClient(c, http.StatusBadRequest, "Request params not valid", err.Error())
		return
	}
	vehicleUser, err, statusError := controller.Service.AttachVehicle(&request, c)
	if err != nil {
		log.Printf("AttachVehicle got error: %s\n", err)
		res.ResponseClient(c, statusError, "AttachVehicle failed", err.Error())
		return
	}
	res.ResponseClient(c, 200, "Attach Vehicle successful", vehicleUser)
}

func (controller *Controller) GetListVehicle(c *gin.Context) {
	vehicles, err, statusError := controller.Service.GetListVehicle(c)
	if err != nil {
		log.Printf("GetListVehicle got error: %s\n", err)
		res.ResponseClient(c, statusError, "GetListVehicle failed", err.Error())
		return
	}
	res.ResponseClient(c, 200, "Get List Vehicle successful", vehicles)
}