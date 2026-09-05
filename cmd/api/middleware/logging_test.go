package middleware_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fiap/postech-tc1/cmd/api/middleware"
	"github.com/gin-gonic/gin"
)

func newRouter(buf *bytes.Buffer, skip ...string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	base := slog.New(slog.NewJSONHandler(buf, nil))

	router := gin.New()
	router.Use(middleware.CorrelationID())
	router.Use(middleware.RequestLogger(base, skip...))

	return router
}

func TestCorrelationID_GeneratesWhenAbsent(t *testing.T) {
	buf := &bytes.Buffer{}
	router := newRouter(buf)

	var seen string
	router.GET("/ping", func(c *gin.Context) {
		seen = middleware.CorrelationIDFrom(c)
		c.Status(http.StatusOK)
	})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/ping", nil))

	if seen == "" {
		t.Fatal("esperava um correlation id gerado")
	}
	if header := rec.Header().Get(middleware.CorrelationHeader); header != seen {
		t.Errorf("esperava o header %q, veio %q", seen, header)
	}
}

func TestCorrelationID_ReusesIncomingHeader(t *testing.T) {
	buf := &bytes.Buffer{}
	router := newRouter(buf)

	router.GET("/ping", func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set(middleware.CorrelationHeader, "trace-externo")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if got := rec.Header().Get(middleware.CorrelationHeader); got != "trace-externo" {
		t.Errorf("esperava %q, veio %q", "trace-externo", got)
	}

	var entry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("log nao e json valido: %v", err)
	}
	if entry["correlation_id"] != "trace-externo" {
		t.Errorf("esperava correlation_id no log, veio %v", entry["correlation_id"])
	}
}

func TestRequestLogger_LevelByStatus(t *testing.T) {
	cases := []struct {
		path   string
		status int
		level  string
	}{
		{"/ok", http.StatusOK, "INFO"},
		{"/nao-encontrado", http.StatusNotFound, "WARN"},
		{"/quebrado", http.StatusInternalServerError, "ERROR"},
	}

	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			buf := &bytes.Buffer{}
			router := newRouter(buf)
			router.GET(tc.path, func(c *gin.Context) { c.Status(tc.status) })

			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tc.path, nil))

			var entry map[string]any
			if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
				t.Fatalf("log nao e json valido: %v", err)
			}
			if entry["level"] != tc.level {
				t.Errorf("esperava level %q, veio %v", tc.level, entry["level"])
			}
			if entry["path"] != tc.path {
				t.Errorf("esperava path %q, veio %v", tc.path, entry["path"])
			}
			if _, ok := entry["latency_ms"]; !ok {
				t.Error("esperava latency_ms no log")
			}
		})
	}
}

func TestRequestLogger_SkipsConfiguredPaths(t *testing.T) {
	buf := &bytes.Buffer{}
	router := newRouter(buf, "/health")
	router.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))

	if buf.Len() != 0 {
		t.Errorf("esperava nenhum log para /health, veio %q", buf.String())
	}
}

func TestLoggerFrom_FallsBackToDefault(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	if middleware.LoggerFrom(c) == nil {
		t.Fatal("esperava um logger, veio nil")
	}
}
