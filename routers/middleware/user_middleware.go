package middleware

import (
	"maintenance-system-go/database/redis"

	"github.com/gin-gonic/gin"
)

type UserMiddleware struct {
	Redis *redis.RedisClient
}
func InitUserMiddleware(redisClient *redis.RedisClient) *UserMiddleware {
	return &UserMiddleware{
		Redis: redisClient,
	}
}
func (um *UserMiddleware) HasLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Looking for cookies in request
		cookie, err := c.Cookie("session_id")
		if err != nil || cookie == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		//Looking user and set current_user for later service
		user, err := um.Redis.Get(c,cookie)
		if err != nil || user == nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		c.Set("current_user", user)

		c.Next()
	}
}