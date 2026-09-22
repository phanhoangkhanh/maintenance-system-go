package user

import (
	"maintenance-system-go/config"

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

func InitUser(db *gorm.DB, config *config.Config) *User {
	repo := &Repository{
		db: db,
	}
	service := &Service{
		Repo: repo,
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
