package app

import (
	"maintenance-system-go/config"
	database "maintenance-system-go/database/connect"
	"maintenance-system-go/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	Database *gorm.DB
	Config   *config.Config
	User     *user.User
	Router   *gin.Engine
}

func InitApp() *App {

	config := config.LoadConfig()
	database := database.ConnectToDatabase(config)
	ginRouter := gin.Default()

	// Register all Modules and Dependencies
	user := user.InitUser(database, config)

	return &App{
		Database: database,
		Config:   config,
		User:     user,
		Router:   ginRouter,
	}
}

// implement indirect.AppContainer interface
func (app *App) ReturnItself() *App {
	return app
}
