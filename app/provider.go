package app

import (
	"maintenance-system-go/config"
	database "maintenance-system-go/database/connect"
	myredis "maintenance-system-go/database/redis"
	"maintenance-system-go/user"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type App struct {
	Redis    *redis.Client
	Database *gorm.DB
	Config   *config.Config
	User     *user.User
	Router   *gin.Engine
}

func InitApp() *App {

	config := config.LoadConfig()
	database := database.ConnectToDatabase(config)
	redisClient := myredis.InitRedis(config)
	ginRouter := gin.Default()

	// Register all Modules and Dependencies
	user := user.InitUser(database, config)

	return &App{
		Redis:    redisClient,
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
