package routers

import (
	"context"
	"maintenance-system-go/config"
	"net/http"
	"slices"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Router struct {
	config *config.Config
	Engine *gin.Engine
}
func InitRouter(config *config.Config) *Router {

	ginRouter := gin.Default()
	//Context for Time out ->pass to all flow
	ginRouter.Use(timeoutMiddleware(2 * time.Second)) 
	//Handle CORS for UI calling
	ginRouter.Use(corsMiddleware(config.CORS))
	// set SameSite policy for all cookies gin context
	ginRouter.Use(func(c *gin.Context) {
		// Default SameSite policy for cookies set during this request.
		c.SetSameSite(http.SameSiteLaxMode)
		c.Next()
	})

	router := &Router{
		config: config,
		Engine: ginRouter,
	}
	return router
}

func corsMiddleware(cfg config.CORSConfig) gin.HandlerFunc {
	options := cors.DefaultConfig()
	// Exact matching also prevents a configured "*" from allowing all origins.
	options.AllowOriginFunc = func(origin string) bool {
		return slices.Contains(cfg.AllowedOrigins, origin)
	}
	options.AllowCredentials = true
	options.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	return cors.New(options)
}

func timeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
  return func(c *gin.Context) {
    ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
    defer cancel()

    // Replace the request with one that carries the new context.
    c.Request = c.Request.WithContext(ctx)
    c.Next()
  }
}
