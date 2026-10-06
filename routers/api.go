package routers

import (
	"maintenance-system-go/app"
	"maintenance-system-go/controller"
)

func (r *Router) RegisterAPIRoutes(app *app.App) {
	r.Engine.GET("/ping", controller.HandleTest)

	// USER APIs
	userController := app.User.Controller
	r.Engine.POST("/login", userController.Login)

	//After Login
	hasLoginRoute := r.Engine.Group("/v1") 
	userMiddleware := app.Middleware.User
	hasLoginRoute.Use(userMiddleware.HasLogin())
	{
		//User
		hasLoginRoute.GET("/user", userController.GetListUser)
		hasLoginRoute.POST("/user", userMiddleware.HasRoleToHandle([]string{"admin", "manager"}), userController.CreateOrUpdateUser)

		//Vehicle
		vehicleController := app.Vehicle.Controller
		//hasLoginRoute.GET("/vehicle", vehicleController.GetListVehicle)
		hasLoginRoute.POST("/vehicle", userMiddleware.HasRoleToHandle([]string{"admin", "technician"}), vehicleController.CreateOrUpdateVehicle)
		hasLoginRoute.POST("/attach-vehicle", userMiddleware.HasRoleToHandle([]string{"admin", "manager"}), vehicleController.AttachVehicle)
	}
}
