package serviceorderhandler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
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

func setupHandler(ctrl *gomock.Controller) (
	*ServiceOrderHandler,
	*mocks.MockCreateServiceOrderUseCase,
	*mocks.MockGetServiceOrderUseCase,
	*mocks.MockListServiceOrdersUseCase,
	*mocks.MockUpdateServiceOrderStatusUseCase,
	*mocks.MockUpdateServiceOrderUseCase,
	*mocks.MockDeleteServiceOrderUseCase,
) {
	create := mocks.NewMockCreateServiceOrderUseCase(ctrl)
	getByID := mocks.NewMockGetServiceOrderUseCase(ctrl)
	getByCode := mocks.NewMockGetServiceOrderByCodeUseCase(ctrl)
	listAll := mocks.NewMockListServiceOrdersUseCase(ctrl)
	listByCPF := mocks.NewMockListServiceOrdersByCPFUseCase(ctrl)
	updateStatus := mocks.NewMockUpdateServiceOrderStatusUseCase(ctrl)
	update := mocks.NewMockUpdateServiceOrderUseCase(ctrl)
	del := mocks.NewMockDeleteServiceOrderUseCase(ctrl)
	approve := mocks.NewMockApproveServiceOrderUseCase(ctrl)
	reject := mocks.NewMockRejectServiceOrderUseCase(ctrl)
	h := NewServiceOrderHandler(create, getByID, getByCode, listAll, listByCPF, updateStatus, update, del, approve, reject)
	return h, create, getByID, listAll, updateStatus, update, del
}

// --- Create ---

func TestServiceOrderHandler_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, create, _, _, _, _, _ := setupHandler(ctrl)
	so := entities.NewServiceOrder("cust-id", "veh-id")
	create.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(so, nil)

	body, _ := json.Marshal(map[string]any{
		"customer_cpf":  "52998224725",
		"vehicle_plate": "ABC1234",
		"notes":         "trocar pastilhas",
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/service-orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestServiceOrderHandler_Create_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, _, _, _, _ := setupHandler(ctrl)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/service-orders", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestServiceOrderHandler_Create_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, create, _, _, _, _, _ := setupHandler(ctrl)
	create.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, domainerrors.ErrNotFound)

	body, _ := json.Marshal(map[string]any{
		"customer_cpf":  "52998224725",
		"vehicle_plate": "ABC1234",
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/service-orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// --- FindByID ---

func TestServiceOrderHandler_FindByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	so := entities.NewServiceOrder("cust-id", "veh-id")
	h, _, getByID, _, _, _, _ := setupHandler(ctrl)
	getByID.EXPECT().Execute(gomock.Any(), "so-1").Return(so, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/service-orders/so-1", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestServiceOrderHandler_FindByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, getByID, _, _, _, _ := setupHandler(ctrl)
	getByID.EXPECT().Execute(gomock.Any(), "so-x").Return(nil, domainerrors.ErrNotFound)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/service-orders/so-x", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// --- FindAll ---

func TestServiceOrderHandler_FindAll_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orders := []*entities.ServiceOrder{entities.NewServiceOrder("c1", "v1")}
	h, _, _, listAll, _, _, _ := setupHandler(ctrl)
	listAll.EXPECT().Execute(gomock.Any()).Return(orders, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/service-orders", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestServiceOrderHandler_FindAll_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, listAll, _, _, _ := setupHandler(ctrl)
	listAll.EXPECT().Execute(gomock.Any()).Return(nil, errors.New("db error"))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/service-orders", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// --- UpdateStatus ---

func TestServiceOrderHandler_UpdateStatus_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, _, updateStatus, _, _ := setupHandler(ctrl)
	updateStatus.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil)

	body, _ := json.Marshal(map[string]any{"status": "in_diagnosis"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/service-orders/so-1/status", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d: %s", w.Code, w.Body.String())
	}
}

func TestServiceOrderHandler_UpdateStatus_InvalidTransition(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, _, updateStatus, _, _ := setupHandler(ctrl)
	updateStatus.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(domainerrors.ErrInvalidStatus)

	body, _ := json.Marshal(map[string]any{"status": "delivered"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/service-orders/so-1/status", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}
}

// --- Delete ---

func TestServiceOrderHandler_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, _, _, _, del := setupHandler(ctrl)
	del.EXPECT().Execute(gomock.Any(), "so-1").Return(nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/service-orders/so-1", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func TestServiceOrderHandler_Delete_NotCancellable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, _, _, _, del := setupHandler(ctrl)
	del.EXPECT().Execute(gomock.Any(), "so-1").Return(domainerrors.ErrOrderNotCancellable)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/service-orders/so-1", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}
}

func TestServiceOrderHandler_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, _, _, _, del := setupHandler(ctrl)
	del.EXPECT().Execute(gomock.Any(), "so-x").Return(domainerrors.ErrNotFound)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/service-orders/so-x", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
