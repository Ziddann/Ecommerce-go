package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string
	GinMode string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	JWTExpiresInHours int
}

func LoadConfig() *Config {
	// Load environment variables
	_ = godotenv.Load()

	expiresIn, err := strconv.Atoi(
		getEnv("JWT_EXPIRES_IN_HOURS", "24"),
	)

	if err != nil || expiresIn <= 0 {
    	expiresIn = 24
	}

	return &Config{
		AppPort: getEnv("APP_PORT", "8080"),
		GinMode: getEnv("GIN_MODE", "debug"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "ecommerce"),
		JWTSecret:  getEnv("JWT_SECRET", ""),
		JWTExpiresInHours: expiresIn,
	}

}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
