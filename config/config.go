package config

import (
	"github.com/fiap/postech-tc1/pkg/env"
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
	MaxFailedLogins     int
	LoginLockMin        int
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		AppPort:             env.GetOrDefault("APP_PORT", "8080"),
		AppEnv:              env.GetOrDefault("APP_ENV", "development"),
		PostgresDSN:         env.GetOrDefault("POSTGRES_DSN", "postgres://postgres:postgres@localhost:5432/workshop?sslmode=disable"),
		JWTSecret:           env.GetOrDefault("JWT_SECRET", "secret"),
		AccessTokenExpMin:   env.GetIntOrDefault("ACCESS_TOKEN_EXP_MIN", 15),
		RefreshTokenExpDays: env.GetIntOrDefault("REFRESH_TOKEN_EXP_DAYS", 7),
		BcryptCost:          env.GetIntOrDefault("BCRYPT_COST", 12),
		MaxFailedLogins:     env.GetIntOrDefault("MAX_FAILED_LOGINS", 5),
		LoginLockMin:        env.GetIntOrDefault("LOGIN_LOCK_MIN", 15),
	}
}
