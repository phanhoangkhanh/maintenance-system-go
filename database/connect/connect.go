package database

import (
	"fmt"
	"maintenance-system-go/config"
	"strconv"

	_postgres "gorm.io/driver/postgres"
	_gorm "gorm.io/gorm"
)

func ConnectToDatabase(cfg *config.Config) *_gorm.DB{
	
	dataConnectString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Ho_Chi_Minh",
		cfg.Database.Host,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Database,
		strconv.Itoa(cfg.Database.Port),
		"disable",
	)

	db, err := _gorm.Open(_postgres.Open(dataConnectString), &_gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	return db
}