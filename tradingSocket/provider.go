package tradingSocket

import (
	"maintenance-system-go/config"
	"maintenance-system-go/database/redis"

	"gorm.io/gorm"
)


type Trading struct {
	Service    *Service
	Database *gorm.DB
}


func InitTrading(db *gorm.DB, config *config.Config, redisClient *redis.RedisClient) *Trading {

	service := &Service{
		Redis: redisClient,
		Config:   config,
	}

	return &Trading{
		Database: db,
		Service:    service,
	}
}