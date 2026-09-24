package app

import (
	"fmt"
	"log"
	"maintenance-system-go/config"
	database "maintenance-system-go/database/connect"
	myredis "maintenance-system-go/database/redis"
	"maintenance-system-go/routers/middleware"
	"maintenance-system-go/user"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	Redis    *myredis.RedisClient
	Database *gorm.DB
	Config   *config.Config
	User     *user.User
	Router   *gin.Engine
	MiddlewareUser *middleware.UserMiddleware
}

func InitApp() *App {

	config := config.LoadConfig()
	database := database.ConnectToDatabase(config)
	redisClient := myredis.InitRedis(config)
	ginRouter := gin.Default()
	//Handle CORS for UI calling
	ginRouter.Use(corsMiddleware(config.CORS))
	// set SameSite policy for all cookies gin context
	ginRouter.Use(func(c *gin.Context) {
		// Default SameSite policy for cookies set during this request.
		c.SetSameSite(http.SameSiteLaxMode)
		c.Next()
	})

	// Register all Modules and Dependencies
	user := user.InitUser(database, config, redisClient)
	userMiddleware := middleware.InitUserMiddleware(redisClient)

	return &App{
		Redis:    redisClient,
		Database: database,
		Config:   config,
		User:     user,
		Router:   ginRouter,
		MiddlewareUser: userMiddleware,
	}
}

func (app *App) CloseApp() {
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
