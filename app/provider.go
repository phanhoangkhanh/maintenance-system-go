package app

import (
	"fmt"
	"log"
	"maintenance-system-go/config"
	database "maintenance-system-go/database/connect"
	myredis "maintenance-system-go/database/redis"
	"maintenance-system-go/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	Redis    *myredis.RedisClient
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
	user := user.InitUser(database, config, redisClient)

	return &App{
		Redis:    redisClient,
		Database: database,
		Config:   config,
		User:     user,
		Router:   ginRouter,
	}
}

func (app *App) CloseApp()  {
	var errs []error

	if app.Redis != nil {
		if err := app.Redis.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close redis: %w", err))
		}
	}

	if app.Database != nil {
		sqlDB, err := app.Database.DB()
		if err != nil {
			errs = append(errs, fmt.Errorf("get sql pool: %w", err))
		} else if err := sqlDB.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close database: %w", err))
		}
	}

	if len(errs) > 0 {
		log.Printf("cleanup app: %v", errs)
	}
}

// implement indirect.AppContainer interface
func (app *App) ReturnItself() *App {
	return app
}
