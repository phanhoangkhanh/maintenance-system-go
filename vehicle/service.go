package vehicle

import (
	"context"
	"fmt"
	"maintenance-system-go/database/redis"
	"maintenance-system-go/models"

	"github.com/gin-gonic/gin"
)

type UserRepoInterface interface {
	GetUser(whereStruct models.Query, c context.Context) ([]models.User, error)
} 

type Service struct {
	Repo  *Repository
	Redis *redis.RedisClient 
	UserRepo UserRepoInterface
}



func (service *Service) CreateNewVehicle(request *CreateOrUpdateVehicleRequest, c *gin.Context) (*models.Vehicle, error, int) {
	//VALIDATE LOGIC
	newVehicle := &models.Vehicle{
		Kind: request.Kind,
		NumberVehicle: request.NumberVehicle,
		Status: request.Status,
		Note: request.Note,
	}
	// Call the repository to save the new vehicle
	vehicle, err := service.Repo.CreateVehicle(newVehicle, c.Request.Context())
	if err != nil {
		return nil, err, 500
	}

	return vehicle, nil, 200
}

func (service *Service) AttachVehicle(request *AttachVehicleRequest, c *gin.Context) (*models.VehicleUser, error, int) {
	//VALIDATE LOGIC
	queryVehicle := &models.Query{
		Where :  []models.WhereClause{
			{
				Key     :"id",
				Compare: "=",
				Value:   request.VehicleID,
			},
		},
		Page: 1,
		PerPage: 1,
	}
	vehicle, err := service.Repo.GetVehicle(*queryVehicle, c.Request.Context())
	if err != nil {
		return nil, err, 500
	}
	if len(vehicle) == 0 {
		return nil, fmt.Errorf("vehicle not found"), 404
	}
	vehicleUpdate := vehicle[0]

	queryVehicle = &models.Query{
		Where :  []models.WhereClause{
			{
				Key     :"id",
				Compare: "=",
				Value:   request.OperatorID,
			},
		},
		Page: 1,
		PerPage: 1,
	}
	operator, err := service.UserRepo.GetUser(*queryVehicle, c.Request.Context())
	if err != nil {
		return nil, err, 500
	}
	if len(operator) == 0 {
		return nil, fmt.Errorf("operator not found"), 404
	}
	operatorUpdate := operator[0]

	queryVehicle = &models.Query{
		Where :  []models.WhereClause{
			{
				Key     :"id",
				Compare: "=",
				Value:   request.DriverID,
			},
		},
		Page: 1,
		PerPage: 1,
	}
	driver, err := service.UserRepo.GetUser(*queryVehicle, c.Request.Context())
	if err != nil {
		return nil, err, 500
	}
	if len(driver) == 0 {
		return nil, fmt.Errorf("driver not found"), 404
	}
	driverUpdate := driver[0]

	newVehicleUser := &models.VehicleUser{
		VehicleID:  vehicleUpdate.ID,
		OperatorID: operatorUpdate.ID,
		DriverID:   driverUpdate.ID,
		Note:       request.Note,
	}

	vehicleUser, err := service.Repo.AttachVehicle(newVehicleUser, c.Request.Context())
	if err != nil {
		return nil, err, 500
	}
	return vehicleUser, nil, 200
}