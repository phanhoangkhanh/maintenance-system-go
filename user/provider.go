package user

import (
	"maintenance-system-go/config"
	"maintenance-system-go/database/redis"

	"gorm.io/gorm"
)

//Controller -> service -> Repo

type User struct {
	Database *gorm.DB
	Config   *config.Config

	Controller *Controller
	Service    *Service
	Repository *Repository
}

func InitUser(db *gorm.DB, config *config.Config, redisClient *redis.RedisClient) *User {
	repo := &Repository{
		db: db,
	}
	service := &Service{
		Repo: repo,
		Redis: redisClient,
	}
	controller := &Controller{
		Service: service,
	}

	return &User{
		Database: db,
		Config:   config,

		Controller: controller,
		Service:    service,
		Repository: repo,
	}
}
