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
	
	if os.Getenv("DB_HOST")== ""  || os.Getenv("DB_PORT") == "" || os.Getenv("DB_USERNAME") == "" || os.Getenv("DB_PASSWORD") == "" || os.Getenv("DB_DATABASE") == "" {
		panic("One or more database environment variables are not set")
	}
	
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
	if os.Getenv("REDIS_HOST") == "" || os.Getenv("REDIS_PORT") == "" || os.Getenv("REDIS_USERNAME") == "" || os.Getenv("REDIS_DB") == "" {
		panic("One or more Redis environment variables are not set")
	}
	return RedisConfig{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     func() int { p, _ := strconv.Atoi(os.Getenv("REDIS_PORT")); return p }(),
		Username: os.Getenv("REDIS_USERNAME"),
		DB:       func() int { d, _ := strconv.Atoi(os.Getenv("REDIS_DB")); return d }(),
	}
}


//TRADING WS CONFIGURATION
type TradingWSConfig struct {
	URI string
}

func LoadTradingWSConfig() TradingWSConfig {
	if url := os.Getenv("TRADING_WS_URI"); url == "" {
		panic("TRADING_WS_URI environment variable is not set")
	}
	return TradingWSConfig{
		URI: os.Getenv("TRADING_WS_URI"),
	}
}