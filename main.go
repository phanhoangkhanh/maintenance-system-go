package main

import (
	"maintenance-system-go/app"
	"maintenance-system-go/routers"
)

func main() {
  
  app := app.InitApp()
  defer app.CloseApp()
  // Register API routes
  routers.RegisterAPIRoutes(app)

  app.Router.Run() // listens on 0.0.0.0:8080 by default
}