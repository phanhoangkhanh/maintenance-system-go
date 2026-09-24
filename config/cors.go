package config

import (
	"os"
	"strings"
)

type CORSConfig struct {
	AllowedOrigins []string
}

func LoadCORSConfig() CORSConfig {
	raw := os.Getenv("CORS_ALLOWED_ORIGINS")
	if strings.TrimSpace(raw) == "" && os.Getenv("ENV") != "prod" {
		raw = "http://localhost:3000,http://localhost:5173"
	}
	var origins []string
	for _, value := range strings.Split(raw, ",") {
		if origin := strings.TrimSpace(value); origin != "" {
			origins = append(origins, origin)
		}
	}
	return CORSConfig{AllowedOrigins: origins}
}
