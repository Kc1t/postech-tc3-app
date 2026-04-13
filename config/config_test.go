package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	os.Unsetenv("APP_PORT")
	os.Unsetenv("APP_ENV")
	os.Unsetenv("POSTGRES_DSN")
	os.Unsetenv("JWT_SECRET")

	cfg := Load()

	if cfg.AppPort != "8080" {
		t.Errorf("expected AppPort %q, got %q", "8080", cfg.AppPort)
	}
	if cfg.AppEnv != "development" {
		t.Errorf("expected AppEnv %q, got %q", "development", cfg.AppEnv)
	}
	if cfg.JWTExpirationHours != 24 {
		t.Errorf("expected JWTExpirationHours 24, got %d", cfg.JWTExpirationHours)
	}
	if cfg.PostgresDSN == "" {
		t.Error("expected non-empty PostgresDSN")
	}
}

func TestLoad_EnvOverride(t *testing.T) {
	os.Setenv("APP_PORT", "9090")
	os.Setenv("APP_ENV", "production")
	os.Setenv("JWT_SECRET", "my-secret")
	defer func() {
		os.Unsetenv("APP_PORT")
		os.Unsetenv("APP_ENV")
		os.Unsetenv("JWT_SECRET")
	}()

	cfg := Load()

	if cfg.AppPort != "9090" {
		t.Errorf("expected AppPort %q, got %q", "9090", cfg.AppPort)
	}
	if cfg.AppEnv != "production" {
		t.Errorf("expected AppEnv %q, got %q", "production", cfg.AppEnv)
	}
	if cfg.JWTSecret != "my-secret" {
		t.Errorf("expected JWTSecret %q, got %q", "my-secret", cfg.JWTSecret)
	}
}

func TestGetEnv_WithValue(t *testing.T) {
	os.Setenv("TEST_KEY_XYZ", "test-value")
	defer os.Unsetenv("TEST_KEY_XYZ")

	result := getEnv("TEST_KEY_XYZ", "fallback")
	if result != "test-value" {
		t.Errorf("expected %q, got %q", "test-value", result)
	}
}

func TestGetEnv_WithFallback(t *testing.T) {
	os.Unsetenv("TEST_KEY_MISSING")
	result := getEnv("TEST_KEY_MISSING", "fallback")
	if result != "fallback" {
		t.Errorf("expected %q, got %q", "fallback", result)
	}
}
