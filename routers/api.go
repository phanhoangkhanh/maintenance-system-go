package routers

import (
	"maintenance-system-go/app"
	"maintenance-system-go/controller"
)

func RegisterAPIRoutes(app *app.App) {
	app.Router.GET("/ping", controller.HandleTest)

	// USER APIs
	userController := app.User.Controller
	app.Router.POST("/login", userController.Login)
}
