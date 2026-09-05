package logger_test

import (
	"log/slog"
	"testing"

	"github.com/fiap/postech-tc1/pkg/logger"
)

func TestParseLevel(t *testing.T) {
	cases := map[string]slog.Level{
		"debug":    slog.LevelDebug,
		"DEBUG":    slog.LevelDebug,
		"  warn  ": slog.LevelWarn,
		"warning":  slog.LevelWarn,
		"error":    slog.LevelError,
		"info":     slog.LevelInfo,
		"":         slog.LevelInfo,
		"qualquer": slog.LevelInfo,
	}

	for input, expected := range cases {
		if got := logger.ParseLevel(input); got != expected {
			t.Errorf("ParseLevel(%q): esperava %v, veio %v", input, expected, got)
		}
	}
}

func TestNew(t *testing.T) {
	if logger.New("debug") == nil {
		t.Fatal("esperava um logger, veio nil")
	}
}
