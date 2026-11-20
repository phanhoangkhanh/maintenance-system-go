package main

import (
	"context"
	"log"
	"maintenance-system-go/app"
	"maintenance-system-go/routers"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	//1-Init App
	app := app.InitApp()
	defer app.CloseApp()

	//2-InitRouter
	router := routers.InitRouter(app.Config)
	// Combine router with app
	router.RegisterAPIRoutes(app)

	// listen signal exit (Ctrl+C)
	ctxClosed, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	//3- Run App (http and socket)
	if err := app.RunApp(router.Engine, ctxClosed); err != nil {
		log.Printf("application stopped: %v", err)
	}
	//no need -> run inside RunApp
	//router.Engine.Run() // listens on 0.0.0.0:8080 by default

}
