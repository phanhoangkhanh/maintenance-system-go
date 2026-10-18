package plan

import (
	"context"
	"fmt"
	"maintenance-system-go/models"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type UserRepository interface{
	GetUser(whereStruct models.Query, c context.Context) ([]models.User, error)
}
type CustomerRepository interface{
	GetCustomer(whereStruct models.Query, c context.Context) ([]models.Customer, error)
}
type VehicleRepository interface{
	GetVehicle(whereStruct models.Query, c context.Context) ([]models.Vehicle, error)
}

type Validator struct {
	Repo *Repository
	UserRepo UserRepository
	CustomerRepo CustomerRepository
	VehicleRepo VehicleRepository
}

func RegisterValidatorPlan(repo *Repository, userRepo UserRepository, customerRepo CustomerRepository, vehicleRepo VehicleRepository) error {
	if repo == nil || repo.db == nil {
		return fmt.Errorf("validator requires a repository with a database")
	}
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return fmt.Errorf("unsupported plan validation engine")
	}
	
	planValidator := &Validator{Repo: repo, UserRepo: userRepo, CustomerRepo: customerRepo, VehicleRepo: vehicleRepo}
	v.RegisterValidation("inCustomerList", planValidator.CustomerExists)
	v.RegisterValidation("inNumberVehicleList", planValidator.NumberVehicleExists)
	v.RegisterValidation("inDriverList", planValidator.DriverExists)
	return nil
}

func (v *Validator) CustomerExists(fl validator.FieldLevel) bool {
	customerName, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	whereClause := models.WhereClause{Key: "name", Compare: "=", Value: customerName }
	queryStruct := models.Query{Where: []models.WhereClause{whereClause}}

	customers, err := v.CustomerRepo.GetCustomer(queryStruct, context.Background())
	if err != nil || len(customers) == 0 {
		return false
	}
	return true
}

func (v *Validator) NumberVehicleExists(fl validator.FieldLevel) bool {
	numberVehicle, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	whereClause := models.WhereClause{Key: "number_vehicle", Compare: "=", Value: numberVehicle }
	queryStruct := models.Query{Where: []models.WhereClause{whereClause}}
	vehicles, err := v.VehicleRepo.GetVehicle(queryStruct, context.Background())
	if err != nil || len(vehicles) == 0 {
		return false
	}
	return true
}

func (v *Validator) DriverExists(fl validator.FieldLevel) bool {
	driverCode, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	whereClause := models.WhereClause{Key: "code", Compare: "=", Value: driverCode }
	queryStruct := models.Query{Where: []models.WhereClause{whereClause}}
	drivers, err := v.UserRepo.GetUser(queryStruct, context.Background())
	if err != nil || len(drivers) == 0 {
		return false
	}
	return true
}