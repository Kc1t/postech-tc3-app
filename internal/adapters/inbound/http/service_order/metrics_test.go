package serviceorderhandler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
)

func TestServiceOrderHandler_AverageExecutionTime_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	m.averageExecutionTime.EXPECT().Execute(gomock.Any()).Return(90*time.Minute, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/service-orders/metrics/execution-time", nil)
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var body struct {
		AverageSeconds float64 `json:"average_seconds"`
		AverageHuman   string  `json:"average_human"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("json decode: %v", err)
	}
	if body.AverageSeconds != (90 * time.Minute).Seconds() {
		t.Errorf("average_seconds = %v, esperava %v", body.AverageSeconds, (90 * time.Minute).Seconds())
	}
	if body.AverageHuman != "1h30m0s" {
		t.Errorf("average_human = %q, esperava %q", body.AverageHuman, "1h30m0s")
	}
}

func TestServiceOrderHandler_AverageExecutionTime_Zero(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	m.averageExecutionTime.EXPECT().Execute(gomock.Any()).Return(time.Duration(0), nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/service-orders/metrics/execution-time", nil)
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var body struct {
		AverageSeconds float64 `json:"average_seconds"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	if body.AverageSeconds != 0 {
		t.Errorf("average_seconds = %v, esperava 0", body.AverageSeconds)
	}
}

func TestServiceOrderHandler_AverageExecutionTime_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	m.averageExecutionTime.EXPECT().Execute(gomock.Any()).Return(time.Duration(0), errors.New("db unavailable"))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/service-orders/metrics/execution-time", nil)
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", w.Code, w.Body.String())
	}
}
