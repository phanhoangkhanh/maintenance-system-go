package config

import (
	"os"
	"strconv"
)

//DATABASE CONFIGURATION
type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

func LoadDBConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     func() int { p, _ := strconv.Atoi(os.Getenv("DB_PORT")); return p }(),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_DATABASE"),
	}
}

//REDIS CONFIGURATION
type RedisConfig struct {
	Host     string
	Port     int
	Username string
	DB       int
}

func LoadRedisConfig() RedisConfig {
	return RedisConfig{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     func() int { p, _ := strconv.Atoi(os.Getenv("REDIS_PORT")); return p }(),
		Username: os.Getenv("REDIS_USERNAME"),
		DB:       func() int { d, _ := strconv.Atoi(os.Getenv("REDIS_DB")); return d }(),
	}
}