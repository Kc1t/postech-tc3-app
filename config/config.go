package config

import (
	"os"
)

type Config struct {
	AppPort            string
	AppEnv             string
	PostgresDSN        string
	JWTSecret          string
	JWTExpirationHours int
}

func Load() *Config {
	return &Config{
		AppPort:            getEnv("APP_PORT", "8080"),
		AppEnv:             getEnv("APP_ENV", "development"),
		PostgresDSN:        getEnv("POSTGRES_DSN", "postgres://postgres:postgres@localhost:5432/workshop?sslmode=disable"),
		JWTSecret:          getEnv("JWT_SECRET", "secret"),
		JWTExpirationHours: 24,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
