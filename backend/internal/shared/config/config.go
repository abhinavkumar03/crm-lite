package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName string
	AppPort string
	AppEnv  string

	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string

	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	JWTSecret string

	FrontendURLs []string
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

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "crm_lite"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),

		JWTSecret: getEnv("JWT_SECRET", "change-me-in-production"),

		FrontendURLs: getEnvAsSlice(
			"FRONTEND_URLS",
			[]string{"http://localhost:3000"},
		),
	}
}

func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	parts := strings.Split(value, ",")

	var result []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}

	return result
}
