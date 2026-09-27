package middleware

import "maintenance-system-go/database/redis"

type Middleware struct {
	User *UserMiddleware
}

func InitMiddleware(redisClient *redis.RedisClient) *Middleware {
	user := InitUserMiddleware(redisClient)
	
	return &Middleware{
		User: user,
	}
}