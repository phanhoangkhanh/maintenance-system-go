package routers

import (
	"encoding/json"
	"fmt"

	"maintenance-system-go/app"
	"maintenance-system-go/controller"
	"maintenance-system-go/models"

	"github.com/gin-gonic/gin"
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
		hasLoginRoute.GET("/profile", func(c *gin.Context) {
			example := c.MustGet("current_user").([]byte)
			var user models.User
			 json.Unmarshal(example, &user)
			fmt.Printf("USER : %v\n", user)
			
			c.JSON(200, user)
		})
	}
}
