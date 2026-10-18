package plan

import (
	"maintenance-system-go/config"

	"gorm.io/gorm"
)
type Plan struct {
	Database *gorm.DB
	Config   *config.Config
	Controller *Controller
	Service    *Service
	Repository *Repository
	Validator   *Validator
}

func InitPlan(db *gorm.DB, config *config.Config, userRepo UserRepository, customerRepo CustomerRepository, vehicleRepo VehicleRepository) *Plan {
	repo := &Repository{
		db: db,
	}
	service := &Service{
		Repo: repo,
	}
	controller := &Controller{
		Service: service,
	}
	validator := &Validator{
		Repo: repo,
		UserRepo: userRepo,
		CustomerRepo: customerRepo,
		VehicleRepo: vehicleRepo,
	}
	//Register Validator from module 
	RegisterValidatorPlan(repo, userRepo, customerRepo, vehicleRepo)

	return &Plan{
		Database: db,
		Config:   config,
		Controller: controller,
		Service:    service,
		Repository: repo,
		Validator:  validator,
	}
}