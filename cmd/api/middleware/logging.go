package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"time"

	"github.com/fiap/postech-tc1/pkg/logger"
	"github.com/gin-gonic/gin"
)

const (
	CorrelationHeader  = "X-Correlation-ID"
	correlationCtxKey  = "correlation_id"
	loggerCtxKey       = "logger"
	fallbackHeaderName = "X-Request-ID"
)

func CorrelationID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(CorrelationHeader)
		if id == "" {
			id = c.GetHeader(fallbackHeaderName)
		}
		if id == "" {
			id = newCorrelationID()
		}

		c.Set(correlationCtxKey, id)
		c.Writer.Header().Set(CorrelationHeader, id)
		c.Next()
	}
}

func RequestLogger(base *slog.Logger, skipPaths ...string) gin.HandlerFunc {
	skip := make(map[string]bool, len(skipPaths))
	for _, p := range skipPaths {
		skip[p] = true
	}

	return func(c *gin.Context) {
		requestLogger := base.With(
			slog.String(correlationCtxKey, CorrelationIDFrom(c)),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
		)
		c.Set(loggerCtxKey, requestLogger)
		c.Request = c.Request.WithContext(logger.ContextWith(c.Request.Context(), requestLogger))

		start := time.Now()
		c.Next()

		if skip[c.Request.URL.Path] {
			return
		}

		attrs := []any{
			slog.Int("status", c.Writer.Status()),
			slog.Int64("latency_ms", time.Since(start).Milliseconds()),
			slog.String("client_ip", c.ClientIP()),
		}

		if userID := c.GetString("user_id"); userID != "" {
			attrs = append(attrs, slog.String("user_id", userID))
		}

		if len(c.Errors) > 0 {
			attrs = append(attrs, slog.String("errors", c.Errors.String()))
		}

		switch {
		case c.Writer.Status() >= 500:
			requestLogger.Error("request", attrs...)
		case c.Writer.Status() >= 400:
			requestLogger.Warn("request", attrs...)
		default:
			requestLogger.Info("request", attrs...)
		}
	}
}

func CorrelationIDFrom(c *gin.Context) string {
	return c.GetString(correlationCtxKey)
}

func LoggerFrom(c *gin.Context) *slog.Logger {
	if c.Request == nil {
		return slog.Default()
	}

	return logger.FromContext(c.Request.Context())
}

func newCorrelationID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "unknown"
	}

	return hex.EncodeToString(b)
}
