package customer

import (
	"context"
	"maintenance-system-go/helper"
	"maintenance-system-go/models"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func InitRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (repo *Repository) ExistsByName(name string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.Customer{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

func (repo *Repository) ExistsByAddress(address string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.Customer{}).Where("address = ?", address).Count(&count).Error
	return count > 0, err
}

func (repo *Repository) Create(request *CustomerRequest, ctx context.Context) (*models.Customer, error) {
	customer := &models.Customer{
		Name:    request.Name,
		Address: request.Address,
		Status:  request.Status,
		Note:    request.Note,
	}
	result := gorm.WithResult()
	err := gorm.G[models.Customer](repo.db, result).Create(ctx, customer)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (repo *Repository) GetCustomer(whereStruct models.Query, c context.Context) ([]models.Customer, error) {
	var customers []models.Customer
	query, err := helper.GrandGetAllInfo(repo.db, whereStruct, c, models.Customer{})
	if err != nil {
		return nil, err
	}
	// Execute the query
	if err := query.Find(&customers).Error; err != nil {
		return nil, err
	}
	

	return customers, nil
}