package main

import (
	"maintenance-system-go/app"
	"maintenance-system-go/routers"
)

func main() {
  
  //1-Init App
  app := app.InitApp()
  defer app.CloseApp()

  //2-InitRouter
  router := routers.InitRouter(app.Config)

  // Combine router with app
  router.RegisterAPIRoutes(app)

  //Keep the server running
  router.Engine.Run() // listens on 0.0.0.0:8080 by default
}