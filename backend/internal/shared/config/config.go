package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// defaultJWTSecret is the placeholder secret shipped for local development. It
// must never be used in production; Load() fails fast if it is.
const defaultJWTSecret = "change-me-in-production"

// FeatureFlags toggles optional subsystems at runtime so incomplete or
// environment-specific capabilities can be shipped safely (config over
// hardcoding). Flags are read once at boot from the environment.
type FeatureFlags struct {
	DynamicModules bool
	Automation     bool
	Import         bool
	Export         bool
	GuidedTour     bool
	RBAC           bool
}

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

	CloudinaryCloudName string
	CloudinaryAPIKey    string
	CloudinaryAPISecret string
	CloudinaryFolder    string

	// WhatsApp / notification provider selection. Defaults keep the app fully
	// functional offline via the simulation provider; setting WhatsAppProvider
	// to "meta" and supplying credentials switches to real delivery.
	WhatsAppProvider string
	WhatsAppAPIURL   string
	WhatsAppToken    string
	WhatsAppPhoneID  string

	JWTSecret     string
	JWTExpiration time.Duration

	FrontendURLs []string

	Features FeatureFlags
}

// Load reads configuration from the environment (and an optional .env file)
// exactly once. It is intended to be called a single time at process start and
// the resulting *Config passed down via dependency injection.
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("config: no .env file found, relying on environment")
	}

	cfg := &Config{
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

		CloudinaryCloudName: getEnv("CLOUDINARY_CLOUD_NAME", ""),
		CloudinaryAPIKey:    getEnv("CLOUDINARY_API_KEY", ""),
		CloudinaryAPISecret: getEnv("CLOUDINARY_API_SECRET", ""),
		CloudinaryFolder:    getEnv("CLOUDINARY_FOLDER", "crm-lite"),

		WhatsAppProvider: getEnv("WHATSAPP_PROVIDER", "simulation"),
		WhatsAppAPIURL:   getEnv("WHATSAPP_API_URL", "https://graph.facebook.com/v20.0"),
		WhatsAppToken:    getEnv("WHATSAPP_TOKEN", ""),
		WhatsAppPhoneID:  getEnv("WHATSAPP_PHONE_ID", ""),

		JWTSecret:     getEnv("JWT_SECRET", defaultJWTSecret),
		JWTExpiration: getEnvAsDuration("JWT_EXPIRATION", 24*time.Hour),

		FrontendURLs: getEnvAsSlice(
			"FRONTEND_URLS",
			[]string{"http://localhost:3000"},
		),

		Features: FeatureFlags{
			DynamicModules: getEnvAsBool("FEATURE_DYNAMIC_MODULES", true),
			Automation:     getEnvAsBool("FEATURE_AUTOMATION", true),
			Import:         getEnvAsBool("FEATURE_IMPORT", true),
			Export:         getEnvAsBool("FEATURE_EXPORT", true),
			GuidedTour:     getEnvAsBool("FEATURE_GUIDED_TOUR", true),
			RBAC:           getEnvAsBool("FEATURE_RBAC", true),
		},
	}

	cfg.mustBeProductionSafe()

	return cfg
}

// IsProduction reports whether the app runs in a production environment.
func (c *Config) IsProduction() bool {
	return strings.EqualFold(c.AppEnv, "production")
}

// mustBeProductionSafe refuses to boot with insecure defaults in production.
func (c *Config) mustBeProductionSafe() {
	if !c.IsProduction() {
		return
	}

	if c.JWTSecret == "" || c.JWTSecret == defaultJWTSecret {
		log.Fatal("config: JWT_SECRET must be set to a strong value in production")
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
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

func getEnvAsBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return boolValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}

	return duration
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	parts := strings.Split(value, ",")

	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}

	if len(result) == 0 {
		return defaultValue
	}

	return result
}
