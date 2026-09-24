package routers

import (
	"encoding/json"
	"fmt"
	"maintenance-system-go/app"
	"maintenance-system-go/controller"
	"maintenance-system-go/models"

	"github.com/gin-gonic/gin"
)

func RegisterAPIRoutes(app *app.App) {
	app.Router.GET("/ping", controller.HandleTest)

	// USER APIs
	userController := app.User.Controller
	app.Router.POST("/login", userController.Login)

	hasLoginRoute := app.Router.Group("/") 
	hasLoginRoute.Use(app.MiddlewareUser.HasLogin())
	{
		hasLoginRoute.GET("/profile", func(c *gin.Context) {
			example := c.MustGet("current_user").([]byte)
			var user models.User
			 json.Unmarshal(example, &user)
			fmt.Printf("USER : %v\n", user)
			
			c.JSON(200, user)
		})
	}
}
