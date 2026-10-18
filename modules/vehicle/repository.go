package vehicle

import (
	"context"
	"maintenance-system-go/helper"
	"maintenance-system-go/models"

	"gorm.io/gorm"
)


type Repository struct {
	db *gorm.DB
}

func (repo *Repository) CreateVehicle(vehicle *models.Vehicle, c context.Context) (*models.Vehicle, error) {
	// Use the provided context when creating the vehicle
	result := gorm.WithResult()
	err := gorm.G[models.Vehicle](repo.db, result).Create(c, vehicle)
	if err != nil {
		return nil, err
	}
	return vehicle, nil
}

func (repo *Repository) GetVehicle(whereStruct models.Query, c context.Context) ([]models.Vehicle, error) {
	var vehicles []models.Vehicle
	query, err := helper.GrandGetAllInfo(repo.db, whereStruct, c, models.Vehicle{})
	if err != nil {
		return nil, err
	}
	//Many to many
	query = query.Preload("DriverAttach").Preload("OperatorAttach")
	// Execute the query
	if err := query.Find(&vehicles).Error; err != nil {
		return nil, err
	}
	

	return vehicles, nil
}

func (repo *Repository) GetVehicleWithConditionEager(whereStruct models.Query, c context.Context, driverName string, operatorName string) ([]models.Vehicle, error) {
	var vehicles []models.Vehicle
	query, err := helper.GrandGetAllInfo(repo.db, whereStruct, c, models.Vehicle{})
	if err != nil {
		return nil, err
	}
	if driverName != ""  {
		//ONly keep the rows adapt conditions inside EXISTS subquery
		query = query.Where(`
			EXISTS (
				SELECT 1
				FROM vehicle_user vu
				JOIN users d ON d.id = vu.driver_id
				WHERE vu.vehicle_id = vehicles.id
				AND d.name LIKE ?
			)`, "%"+driverName+"%")
	}
	if operatorName != ""  {
		query = query.Where(`
			EXISTS (
				SELECT 1
				FROM vehicle_user vu
				JOIN users o ON o.id = vu.operator_id
				WHERE vu.vehicle_id = vehicles.id
				AND o.name LIKE ?
			)`, "%"+operatorName+"%")
	}
	//Many to many
	query = query.Preload("DriverAttach").Preload("OperatorAttach")
	// Execute the query
	if err := query.Find(&vehicles).Error; err != nil {
		return nil, err
	}
	return vehicles, nil
}

func (repo *Repository) AttachVehicle(vehicleUser *models.VehicleUser, c context.Context) (*models.VehicleUser, error) {
	result := gorm.WithResult()
	err := gorm.G[models.VehicleUser](repo.db, result).Create(c, vehicleUser)
	if err != nil {
		return nil, err
	}
	return vehicleUser, nil
}