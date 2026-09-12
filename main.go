package main

import (
	"maintenance-system-go/controller"

	"github.com/gin-gonic/gin"
)

func main() {
  router := gin.Default()
  router.GET("/ping", controller.HandleTest)
  // Submit Login Form
  router.POST("/login", controller.Login)


  router.Run() // listens on 0.0.0.0:8080 by default
}