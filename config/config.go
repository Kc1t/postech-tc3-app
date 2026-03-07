package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort             string
	AppEnv              string
	MongoURI            string
	MongoDB             string
	JWTSecret           string
	JWTExpirationHours  int
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		AppPort:            getEnv("APP_PORT", "8080"),
		AppEnv:             getEnv("APP_ENV", "development"),
		MongoURI:           getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:            getEnv("MONGO_DB", "workshop"),
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
