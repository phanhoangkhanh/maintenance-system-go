package app

import (
	"maintenance-system-go/authenticate"
	"maintenance-system-go/config"
	database "maintenance-system-go/database/connect"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	Database *gorm.DB
	Config *config.Config
	Authenticate *authenticate.Authenticate
	Router *gin.Engine

}

func InitApp() *App {

	config := config.LoadConfig()
  	database := database.ConnectToDatabase(config)
	ginRouter := gin.Default()

	// Register all services and dependencies
	authenticate := authenticate.InitAuthenticate(database, config)

	return &App{
		Database: database,
		Config:   config,
		Authenticate: authenticate,
		Router:    ginRouter,

	}
}

//implement indirect.AppContainer interface
func (app *App) ReturnItself() *App {
	return app
}

