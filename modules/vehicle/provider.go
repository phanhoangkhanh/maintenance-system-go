package vehicle

import (
	"maintenance-system-go/config"
	"maintenance-system-go/database/redis"

	"gorm.io/gorm"
)

//Controller -> service -> Repo

type Vehicle struct {
	Database *gorm.DB
	Config   *config.Config

	Controller *Controller
	Service    *Service
	Repository *Repository
}

func InitVehicle(db *gorm.DB, config *config.Config, redisClient *redis.RedisClient, userRepo UserRepoInterface) *Vehicle {
	repo := &Repository{
		db: db,
	}
	service := &Service{
		Repo: repo,
		Redis: redisClient, 
		UserRepo: userRepo, 
	}
	controller := &Controller{
		Service: service,
	}

	//Register Validator from module Vehicle

	return &Vehicle{
		Database: db,
		Config:   config,

		Controller: controller,
		Service:    service,
		Repository: repo,
	}
}


