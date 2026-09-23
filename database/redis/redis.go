package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	config "maintenance-system-go/config"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}
func InitRedis(cfg *config.Config) *RedisClient {
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
	return &RedisClient{Client: client}
}

func (redis *RedisClient) Create(ctx context.Context, value interface{}, ttl time.Duration) (string, error) {
	// 1- create key ID 
	id, err := uuid.NewRandom() 
	if err != nil {
		return "", fmt.Errorf("Redis generate session ID failed: %w", err)
	}
	sessionID := id.String() 

	//2 - values 
	data, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("Redis marshal value failed: %w", err)
	}

	//3- store in redis
	err = redis.Client.Set(ctx, sessionID, data, ttl).Err()
	if err != nil {
		return "", fmt.Errorf("Redis set session failed: %w", err)
	}
	return sessionID, nil
}

func (redis *RedisClient) Get(ctx context.Context, sessionID string) ([]byte, error){
	data, err := redis.Client.Get(ctx, sessionID).Bytes()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (redis *RedisClient) Close() error {
	err := redis.Client.Close()
	if err != nil {
		return err
	}
	return nil
}
