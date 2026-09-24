package app

import (
	"maintenance-system-go/config"
	"slices"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

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
