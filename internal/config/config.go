package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv        string
	ServerPort    string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	JWTSecret     string
	JWTExpiration time.Duration
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	expirationHours := getEnvInt("JWT_EXPIRATION_HOURS", 24)

	cfg := &Config{
		AppEnv:        getEnv("APP_ENV", "development"),
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        getEnv("DB_PORT", "3306"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		JWTExpiration: time.Duration(expirationHours) * time.Hour,
	}

	if cfg.DBHost == "" || cfg.DBUser == "" || cfg.DBName == "" || cfg.JWTSecret == "" {
		return nil, errors.New("DB_HOST, DB_USER, DB_NAME, and JWT_SECRET are required")
	}

	return cfg, nil
}

func (c *Config) DSN(withDatabase bool) string {
	base := fmt.Sprintf("%s:%s@tcp(%s:%s)/", c.DBUser, c.DBPassword, c.DBHost, c.DBPort)
	if withDatabase {
		base += c.DBName
	}

	return base + "?parseTime=true&charset=utf8mb4,utf8&loc=Local&multiStatements=true"
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return value
}
