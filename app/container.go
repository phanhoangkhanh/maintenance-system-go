package app

import (
	"fmt"
	"log"
	"maintenance-system-go/config"
	database "maintenance-system-go/database/connect"
	myredis "maintenance-system-go/database/redis"
	"maintenance-system-go/middleware"
	"maintenance-system-go/modules/customer"
	"maintenance-system-go/modules/plan"
	"maintenance-system-go/modules/user"
	"maintenance-system-go/modules/vehicle"
	"maintenance-system-go/tradingSocket"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)



type App struct {
	Redis    *myredis.RedisClient
	Database *gorm.DB
	Config   *config.Config
	User     *user.User
	Vehicle  *vehicle.Vehicle
	Customer *customer.Customer
	Plan     *plan.Plan
	Router   *gin.Engine
	Middleware *middleware.Middleware
	Trading *tradingSocket.Trading
}

func InitApp() *App {
	//Load all config params
	config := config.LoadConfig()

	database := database.ConnectToDatabase(config)
	redisClient := myredis.InitRedis(config)

	// Register all Modules and Dependencies
	user := user.InitUser(database, config, redisClient)
	vehicle := vehicle.InitVehicle(database, config, redisClient, user.Repository) // inject real struct
	customer := customer.InitCustomer(database, config) // inject real struct
	plan := plan.InitPlan(database, config, user.Repository, customer.Repository, vehicle.Repository) // inject real struct Repo

	middleware := middleware.InitMiddleware(redisClient)

	//Trading socket 
	trading := tradingSocket.InitTrading(database, config, redisClient)
	//Init routine to keep socket connection alive
	go trading.Service.ConnectWebsocket()

	return &App{
		Redis:    redisClient,
		Database: database,
		Config:   config,
		User:     user,
		Vehicle:  vehicle,
		Customer: customer,
		Plan:     plan,
		Middleware: middleware,
		Trading: trading,
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

	
	//Close trading socket connection if exists
	if app.Trading != nil && app.Trading.Service != nil && app.Trading.Service.Conn != nil {
		if err := app.Trading.Service.Conn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close trading socket: %w", err))
		}
	}


	if len(errs) > 0 {
		log.Printf("cleanup app: %v", errs)
		panic("We got error , but push to close App")
	}
}

// implement indirect.AppContainer interface
func (app *App) ReturnItself() *App {
	return app
}
