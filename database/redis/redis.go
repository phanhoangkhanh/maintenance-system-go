package redis

import (
	"context"
	"fmt"
	"log"

	config "maintenance-system-go/config"

	"github.com/redis/go-redis/v9"
)
func InitRedis(cfg *config.Config) *redis.Client {
	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Username: cfg.Redis.Username,
		DB:       cfg.Redis.DB,
	}
	client := redis.NewClient(options)
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %s", err)
	}
	fmt.Println("Ping successlly Redis")
	return client
}