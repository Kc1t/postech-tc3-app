package logger_test

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/fiap/postech-tc1/pkg/logger"
)

func TestContextRoundTrip(t *testing.T) {
	buf := &bytes.Buffer{}
	l := slog.New(slog.NewJSONHandler(buf, nil))

	ctx := logger.ContextWith(context.Background(), l)
	logger.FromContext(ctx).Info("oi")

	if buf.Len() == 0 {
		t.Error("esperava o logger do contexto ser usado")
	}
}

func TestFromContext_FallsBackToDefault(t *testing.T) {
	if logger.FromContext(context.Background()) == nil {
		t.Error("esperava o logger padrao, veio nil")
	}
}

func TestContextWith_IgnoresNil(t *testing.T) {
	ctx := logger.ContextWith(context.Background(), nil)

	if logger.FromContext(ctx) == nil {
		t.Error("esperava o logger padrao quando o informado e nil")
	}
}
