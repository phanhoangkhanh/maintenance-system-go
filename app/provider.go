package app

import (
	"maintenance-system-go/authenticate"
	"maintenance-system-go/config"
	database "maintenance-system-go/database/connect"

	"gorm.io/gorm"
)

type App struct {
	Database *gorm.DB
	Config *config.Config
	Authenticate *authenticate.Authenticate
}

func InitApp() *App {

	config := config.LoadConfig()
  	database := database.ConnectToDatabase(config)
	authenticate := authenticate.InitAuthenticate(database, config)
	
	return &App{
		Database: database,
		Config:   config,
		Authenticate: authenticate,
	}
}

