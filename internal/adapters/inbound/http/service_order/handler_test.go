package serviceorderhandler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestRouter(h *ServiceOrderHandler) *gin.Engine {
	r := gin.New()
	h.SetupRoutes(r.Group(""))
	return r
}

func newHandler(ctrl *gomock.Controller) *ServiceOrderHandler {
	return NewServiceOrderHandler(
		mocks.NewMockCreateServiceOrderUseCase(ctrl),
		mocks.NewMockGetServiceOrderUseCase(ctrl),
		mocks.NewMockListServiceOrdersUseCase(ctrl),
		mocks.NewMockUpdateServiceOrderStatusUseCase(ctrl),
		mocks.NewMockUpdateServiceOrderUseCase(ctrl),
		mocks.NewMockDeleteServiceOrderUseCase(ctrl),
	)
}

func TestServiceOrderHandler_Create_NotImplemented(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/service-orders", nil)
	newTestRouter(newHandler(ctrl)).ServeHTTP(w, req)

	if w.Code != http.StatusNotImplemented {
		t.Fatalf("expected 501, got %d", w.Code)
	}
}

func TestServiceOrderHandler_FindAll_NotImplemented(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/service-orders", nil)
	newTestRouter(newHandler(ctrl)).ServeHTTP(w, req)

	if w.Code != http.StatusNotImplemented {
		t.Fatalf("expected 501, got %d", w.Code)
	}
}

func TestServiceOrderHandler_FindByID_NotImplemented(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/service-orders/so-1", nil)
	newTestRouter(newHandler(ctrl)).ServeHTTP(w, req)

	if w.Code != http.StatusNotImplemented {
		t.Fatalf("expected 501, got %d", w.Code)
	}
}

func TestServiceOrderHandler_Update_NotImplemented(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/service-orders/so-1", nil)
	newTestRouter(newHandler(ctrl)).ServeHTTP(w, req)

	if w.Code != http.StatusNotImplemented {
		t.Fatalf("expected 501, got %d", w.Code)
	}
}

func TestServiceOrderHandler_UpdateStatus_NotImplemented(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/service-orders/so-1/status", nil)
	newTestRouter(newHandler(ctrl)).ServeHTTP(w, req)

	if w.Code != http.StatusNotImplemented {
		t.Fatalf("expected 501, got %d", w.Code)
	}
}

func TestServiceOrderHandler_Delete_NotImplemented(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/service-orders/so-1", nil)
	newTestRouter(newHandler(ctrl)).ServeHTTP(w, req)

	if w.Code != http.StatusNotImplemented {
		t.Fatalf("expected 501, got %d", w.Code)
	}
}
