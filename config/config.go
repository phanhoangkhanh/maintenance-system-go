package config

import (
	"log"
	"os"

	_godotenv "github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
}


func LoadConfig() *Config {
	if os.Getenv("ENV") != "prod" {
		err := _godotenv.Load(".env")
		if err != nil  {
			log.Fatalf("Error loading .env file %s", err)
		}
	}

	return &Config{
		Database:  LoadDBConfig(),
		
	}
}