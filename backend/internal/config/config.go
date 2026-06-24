package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName string
	AppPort string
	AppEnv  string
}

func Load() *Config {
	err := godotenv.Load()

	if err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		AppName: getEnv("APP_NAME", "CRM Lite"),
		AppPort: getEnv("APP_PORT", "8080"),
		AppEnv:  getEnv("APP_ENV", "development"),
	}
}
