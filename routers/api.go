package routers

import (
	"maintenance-system-go/app"
	"maintenance-system-go/controller"
)	

func  RegisterAPIRoutes(app *app.App) {
	app.Router.GET("/ping", controller.HandleTest)

  	// AUTHENTICATE APIs
	authenticate := app.Authenticate
  	app.Router.POST("/login", authenticate.Login)
} 