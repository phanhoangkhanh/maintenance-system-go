package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	broadcastsocket "maintenance-system-go/broadcastSocket"
	"maintenance-system-go/config"
	database "maintenance-system-go/database/connect"
	myredis "maintenance-system-go/database/redis"
	"maintenance-system-go/middleware"
	"maintenance-system-go/modules/customer"
	"maintenance-system-go/modules/plan"
	"maintenance-system-go/modules/user"
	"maintenance-system-go/modules/vehicle"
	"maintenance-system-go/tradingSocket"
	"net/http"
	"time"

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
	HTTPServer *http.Server // for router init
	Middleware *middleware.Middleware
	Trading *tradingSocket.Trading
	BroadcastSocket *broadcastsocket.BroadCastSocket
}

//Initial Bridge Channel for 2 Sockets 
//This Channel be used to transfer trading-socket info to broadcast-socket to let FE can access info
var BridgeChannel = make(chan tradingSocket.MarketPriceInfo, 1000) 

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
	trading := tradingSocket.InitTrading(database, config, redisClient, BridgeChannel)
	
	//Broadcast socket
	broadcast := broadcastsocket.InitBroadCast(config, BridgeChannel)
	

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
		BroadcastSocket: broadcast,
	}
}

func (app *App) RunApp(engine *gin.Engine, ctxclosed context.Context) error{


	//Run Socket keep listen trading socket
	go app.Trading.Service.ConnectWebsocket() 
	// Run socket to broadcast trading info to clients
	go app.BroadcastSocket.Service.CreatAndBroadcastSocket()

	//Keep wait ctxClosed
	go func(){
		for{
			select{
				case <-ctxclosed.Done():
					log.Println("Context closed, shutting down...")
					app.CloseApp()
					return
				default:
			}
		}
	}()

	// Run HTTP server
	app.HTTPServer = &http.Server{
		Addr:              ":8080",
        Handler:           engine,
        ReadHeaderTimeout: 5 * time.Second,
	}
	//Keep main goroutine running to serve HTTP and sockets
	err := app.HTTPServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
        return fmt.Errorf("serve HTTP: %w", err)
    }

    return nil
}

func (app *App) CloseApp() {
	var errs []error

	if app.Redis != nil {
		if err := app.Redis.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close redis: %w", err))
		}
		log.Println("Close redis connection")
	}

	if app.Database != nil {
		sqlDB, err := app.Database.DB()
		if err != nil {
			errs = append(errs, fmt.Errorf("get sql pool: %w", err))
		} else if err := sqlDB.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close database: %w", err))
		}
		log.Println("Close database connection")
	}

	
	//Close trading socket connection if exists
	if app.Trading != nil && app.Trading.Service != nil && app.Trading.Service.Conn != nil {
		if err := app.Trading.Service.Conn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close trading socket: %w", err))
		}
		log.Println("Close trading socket connection")
	}

	//Close broadcast socket connection if exists
	if app.BroadcastSocket != nil && app.BroadcastSocket.Service != nil && app.BroadcastSocket.Service.Conn != nil {
		if err := app.BroadcastSocket.Service.Conn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close broadcast socket: %w", err))
		}
		log.Println("Close broadcast socket connection")
	}
	// close BridgeChannel if exists and still open
	if BridgeChannel != nil {
		close(BridgeChannel)
		log.Println("Close BridgeChannel")
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
