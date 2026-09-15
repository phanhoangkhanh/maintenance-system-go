package main

import (
	"maintenance-system-go/app"
	"maintenance-system-go/controller"

	"github.com/gin-gonic/gin"
)

func main() {
  
  app := app.InitApp()
  router := gin.Default()
  
  router.GET("/ping", controller.HandleTest)
  // Submit Login Form
  router.POST("/login", app.Authenticate.Login)


  router.Run() // listens on 0.0.0.0:8080 by default
}