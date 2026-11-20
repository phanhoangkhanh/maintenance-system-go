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
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	Redis           *myredis.RedisClient
	Database        *gorm.DB
	Config          *config.Config
	User            *user.User
	Vehicle         *vehicle.Vehicle
	Customer        *customer.Customer
	Plan            *plan.Plan
	HTTPServer      *http.Server // for router init
	Middleware      *middleware.Middleware
	Trading         *tradingSocket.Trading
	BroadcastSocket *broadcastsocket.BroadCastSocket
}

// Initial Bridge Channel for 2 Sockets
// This Channel be used to transfer trading-socket info to broadcast-socket to let FE can access info
var BridgeChannel = make(chan tradingSocket.MarketPriceInfo, 1000)

func InitApp() *App {
	//Load all config params
	config := config.LoadConfig()

	database := database.ConnectToDatabase(config)
	redisClient := myredis.InitRedis(config)

	// Register all Modules and Dependencies
	user := user.InitUser(database, config, redisClient)
	vehicle := vehicle.InitVehicle(database, config, redisClient, user.Repository)                    // inject real struct
	customer := customer.InitCustomer(database, config)                                               // inject real struct
	plan := plan.InitPlan(database, config, user.Repository, customer.Repository, vehicle.Repository) // inject real struct Repo

	middleware := middleware.InitMiddleware(redisClient)

	//Trading socket
	trading := tradingSocket.InitTrading(database, config, redisClient, BridgeChannel)

	//Broadcast socket
	broadcast := broadcastsocket.InitBroadCast(config, BridgeChannel)

	return &App{
		Redis:           redisClient,
		Database:        database,
		Config:          config,
		User:            user,
		Vehicle:         vehicle,
		Customer:        customer,
		Plan:            plan,
		Middleware:      middleware,
		Trading:         trading,
		BroadcastSocket: broadcast,
	}
}

func (app *App) RunApp(engine *gin.Engine, ctxclosed context.Context) error {
	if ctxclosed.Err() != nil {
		return nil
	}

	app.HTTPServer = &http.Server{
		Addr:              ":8080",
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
	}
	// Bind before starting background work, so a busy port fails immediately.
	listener, err := net.Listen("tcp", app.HTTPServer.Addr)
	if err != nil {
		return fmt.Errorf("listen HTTP: %w", err)
	}

	go app.Trading.Service.ConnectWebsocket()
	go app.BroadcastSocket.Service.CreatAndBroadcastSocket()

	// CloseApp runs once through main's defer, after HTTP has stopped.
	return serveHTTP(ctxclosed, app.HTTPServer, listener)
}

func serveHTTP(ctx context.Context, server *http.Server, listener net.Listener) error {
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.Serve(listener)
	}()

	select {
	case err := <-serverErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("serve HTTP: %w", err)
		}
		return nil
	case <-ctx.Done():
		log.Println("Shutdown requested, waiting for HTTP requests to finish...")
	}

	// The signal context is already canceled; shutdown needs its own deadline.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var shutdownErr error
	if err := server.Shutdown(shutdownCtx); err != nil {
		shutdownErr = fmt.Errorf("shutdown HTTP: %w", err)
		if err := server.Close(); err != nil {
			shutdownErr = errors.Join(shutdownErr, fmt.Errorf("close HTTP: %w", err))
		}
	}
	// Wait for the serving goroutine even when shutdown required a forced close.
	if err := <-serverErr; err != nil && !errors.Is(err, http.ErrServerClosed) {
		shutdownErr = errors.Join(shutdownErr, fmt.Errorf("serve HTTP: %w", err))
	}

	return shutdownErr
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
	// Do not close BridgeChannel until all trading senders have stopped.
	// Trading currently reconnects independently, so closing it here can panic.

	if len(errs) > 0 {
		log.Printf("cleanup app: %v", errs)
		panic("We got error , but push to close App")
	}
}

// implement indirect.AppContainer interface
func (app *App) ReturnItself() *App {
	return app
}
