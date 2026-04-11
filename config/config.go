package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort             string
	AppEnv              string
	PostgresDSN         string
	JWTSecret           string
	AccessTokenExpMin   int
	RefreshTokenExpDays int
	BcryptCost          int
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		AppPort:             getEnv("APP_PORT", "8080"),
		AppEnv:              getEnv("APP_ENV", "development"),
		PostgresDSN:         getEnv("POSTGRES_DSN", "postgres://postgres:postgres@localhost:5432/workshop?sslmode=disable"),
		JWTSecret:           getEnv("JWT_SECRET", "secret"),
		AccessTokenExpMin:   getEnvInt("ACCESS_TOKEN_EXP_MIN", 15),
		RefreshTokenExpDays: getEnvInt("REFRESH_TOKEN_EXP_DAYS", 7),
		BcryptCost:          getEnvInt("BCRYPT_COST", 12),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
