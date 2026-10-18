package customer

import (
	"maintenance-system-go/config"

	"gorm.io/gorm"
)

type Customer struct {
	Database *gorm.DB
	Config   *config.Config
	Controller *Controller
	Service    *Service
	Repository *Repository
	Validator   *ValidatorCustomer
}

func InitCustomer(db *gorm.DB, config *config.Config) *Customer {
	repo := &Repository{
		db: db,
	}
	service := &Service{
		Repo: repo,
	}
	controller := &Controller{
		Service: service,
	}
	validator := &ValidatorCustomer{
		Repo: repo,
	}

	//Register Validator from module 
	RegisterValidatorCustomer(repo)

	return &Customer{
		Database: db,
		Config:   config,
		Controller: controller,
		Service:    service,
		Repository: repo,
		Validator:  validator,
	}
}