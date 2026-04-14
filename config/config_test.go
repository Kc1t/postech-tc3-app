package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	t.Helper()
	for _, key := range []string{"APP_PORT", "APP_ENV", "POSTGRES_DSN", "JWT_SECRET"} {
		if err := os.Unsetenv(key); err != nil {
			t.Fatalf("Unsetenv(%q): %v", key, err)
		}
	}

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
	envs := map[string]string{
		"APP_PORT":   "9090",
		"APP_ENV":    "production",
		"JWT_SECRET": "my-secret",
	}
	for k, v := range envs {
		if err := os.Setenv(k, v); err != nil {
			t.Fatalf("Setenv(%q): %v", k, err)
		}
	}
	defer func() {
		for k := range envs {
			_ = os.Unsetenv(k)
		}
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
	if err := os.Setenv("TEST_KEY_XYZ", "test-value"); err != nil {
		t.Fatalf("Setenv: %v", err)
	}
	defer func() { _ = os.Unsetenv("TEST_KEY_XYZ") }()

	result := getEnv("TEST_KEY_XYZ", "fallback")
	if result != "test-value" {
		t.Errorf("expected %q, got %q", "test-value", result)
	}
}

func TestGetEnv_WithFallback(t *testing.T) {
	if err := os.Unsetenv("TEST_KEY_MISSING"); err != nil {
		t.Fatalf("Unsetenv: %v", err)
	}
	result := getEnv("TEST_KEY_MISSING", "fallback")
	if result != "fallback" {
		t.Errorf("expected %q, got %q", "fallback", result)
	}
}
