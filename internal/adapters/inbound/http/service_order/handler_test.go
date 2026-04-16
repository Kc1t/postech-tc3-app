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
	h.SetupPublicRoutes(r.Group(""))
	return r
}

type handlerMocks struct {
	handler      *ServiceOrderHandler
	create       *mocks.MockCreateServiceOrderUseCase
	getByID      *mocks.MockGetServiceOrderUseCase
	getByCode    *mocks.MockGetServiceOrderByCodeUseCase
	listAll      *mocks.MockListServiceOrdersUseCase
	listByCPF    *mocks.MockListServiceOrdersByCPFUseCase
	updateStatus *mocks.MockUpdateServiceOrderStatusUseCase
	update       *mocks.MockUpdateServiceOrderUseCase
	del          *mocks.MockDeleteServiceOrderUseCase
	approve      *mocks.MockApproveServiceOrderUseCase
	reject       *mocks.MockRejectServiceOrderUseCase
}

func newHandler(ctrl *gomock.Controller) handlerMocks {
	m := handlerMocks{
		create:       mocks.NewMockCreateServiceOrderUseCase(ctrl),
		getByID:      mocks.NewMockGetServiceOrderUseCase(ctrl),
		getByCode:    mocks.NewMockGetServiceOrderByCodeUseCase(ctrl),
		listAll:      mocks.NewMockListServiceOrdersUseCase(ctrl),
		listByCPF:    mocks.NewMockListServiceOrdersByCPFUseCase(ctrl),
		updateStatus: mocks.NewMockUpdateServiceOrderStatusUseCase(ctrl),
		update:       mocks.NewMockUpdateServiceOrderUseCase(ctrl),
		del:          mocks.NewMockDeleteServiceOrderUseCase(ctrl),
		approve:      mocks.NewMockApproveServiceOrderUseCase(ctrl),
		reject:       mocks.NewMockRejectServiceOrderUseCase(ctrl),
	}
	m.handler = NewServiceOrderHandler(
		m.create, m.getByID, m.getByCode, m.listAll, m.listByCPF,
		m.updateStatus, m.update, m.del, m.approve, m.reject,
	)
	return m
}

// =============================================================================
// Create
// =============================================================================

func TestServiceOrderHandler_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	so := entities.NewServiceOrder("cust-id", "veh-id")
	m.create.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(so, nil)

	body, _ := json.Marshal(map[string]any{
		"customer_cpf":  "52998224725",
		"vehicle_plate": "ABC1234",
		"notes":         "trocar pastilhas",
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/service-orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestServiceOrderHandler_Create_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/service-orders", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestServiceOrderHandler_Create_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	m.create.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, domainerrors.ErrNotFound)

	body, _ := json.Marshal(map[string]any{
		"customer_cpf":  "52998224725",
		"vehicle_plate": "ABC1234",
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/service-orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestServiceOrderHandler_Create_VehicleNotFromCustomer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	m.create.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, domainerrors.ErrVehicleNotFromCustomer)

	body, _ := json.Marshal(map[string]any{
		"customer_cpf":  "52998224725",
		"vehicle_plate": "ABC1234",
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/service-orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}
}

func TestServiceOrderHandler_Create_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	m.create.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, errors.New("db error"))

	body, _ := json.Marshal(map[string]any{
		"customer_cpf":  "52998224725",
		"vehicle_plate": "ABC1234",
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/service-orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// =============================================================================
// FindByID
// =============================================================================

func TestServiceOrderHandler_FindByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	so := entities.NewServiceOrder("cust-id", "veh-id")
	m.getByID.EXPECT().Execute(gomock.Any(), "so-1").Return(so, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/service-orders/so-1", nil)
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestServiceOrderHandler_FindByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	m.getByID.EXPECT().Execute(gomock.Any(), "so-x").Return(nil, domainerrors.ErrNotFound)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/service-orders/so-x", nil)
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// =============================================================================
// FindAll
// =============================================================================

func TestServiceOrderHandler_FindAll_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	orders := []*entities.ServiceOrder{entities.NewServiceOrder("c1", "v1")}
	m.listAll.EXPECT().Execute(gomock.Any()).Return(orders, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/service-orders", nil)
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestServiceOrderHandler_FindAll_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	m.listAll.EXPECT().Execute(gomock.Any()).Return(nil, errors.New("db error"))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/service-orders", nil)
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// =============================================================================
// UpdateStatus
// =============================================================================

func TestServiceOrderHandler_UpdateStatus_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	m.updateStatus.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil)

	body, _ := json.Marshal(map[string]any{"status": "in_diagnosis"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/service-orders/so-1/status", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d: %s", w.Code, w.Body.String())
	}
}

func TestServiceOrderHandler_UpdateStatus_InvalidTransition(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	m.updateStatus.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(domainerrors.ErrInvalidStatus)

	body, _ := json.Marshal(map[string]any{"status": "delivered"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/service-orders/so-1/status", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}
}

func TestServiceOrderHandler_UpdateStatus_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	m.updateStatus.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(domainerrors.ErrNotFound)

	body, _ := json.Marshal(map[string]any{"status": "in_diagnosis"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/service-orders/so-1/status", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestServiceOrderHandler_UpdateStatus_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/service-orders/so-1/status", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

// =============================================================================
// Update
// =============================================================================

func TestServiceOrderHandler_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	so := entities.NewServiceOrder("cust-id", "veh-id")
	notes := "nova observacao"

	m.getByID.EXPECT().Execute(gomock.Any(), "so-1").Return(so, nil)
	m.update.EXPECT().Execute(gomock.Any(), so).Return(nil)

	body, _ := json.Marshal(map[string]any{"notes": notes})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/service-orders/so-1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestServiceOrderHandler_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	m.getByID.EXPECT().Execute(gomock.Any(), "so-x").Return(nil, domainerrors.ErrNotFound)

	body, _ := json.Marshal(map[string]any{"notes": "x"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/service-orders/so-x", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestServiceOrderHandler_Update_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	so := entities.NewServiceOrder("cust-id", "veh-id")

	m.getByID.EXPECT().Execute(gomock.Any(), "so-1").Return(so, nil)
	m.update.EXPECT().Execute(gomock.Any(), so).Return(errors.New("db error"))

	body, _ := json.Marshal(map[string]any{"notes": "x"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/service-orders/so-1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// =============================================================================
// Delete
// =============================================================================

func TestServiceOrderHandler_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	m.del.EXPECT().Execute(gomock.Any(), "so-1").Return(nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/service-orders/so-1", nil)
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func TestServiceOrderHandler_Delete_NotCancellable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	m.del.EXPECT().Execute(gomock.Any(), "so-1").Return(domainerrors.ErrOrderNotCancellable)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/service-orders/so-1", nil)
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}
}

func TestServiceOrderHandler_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	m.del.EXPECT().Execute(gomock.Any(), "so-x").Return(domainerrors.ErrNotFound)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/service-orders/so-x", nil)
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestServiceOrderHandler_Delete_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	m.del.EXPECT().Execute(gomock.Any(), "so-1").Return(errors.New("db error"))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/service-orders/so-1", nil)
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// =============================================================================
// FindByCode (rota publica)
// =============================================================================

func TestServiceOrderHandler_FindByCode_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	so := entities.NewServiceOrder("cust-id", "veh-id")
	m.getByCode.EXPECT().Execute(gomock.Any(), 100, "52998224725").Return(so, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/service-orders/code/100?cpf=52998224725", nil)
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestServiceOrderHandler_FindByCode_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	m.getByCode.EXPECT().Execute(gomock.Any(), 999, "52998224725").Return(nil, domainerrors.ErrNotFound)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/service-orders/code/999?cpf=52998224725", nil)
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestServiceOrderHandler_FindByCode_SemCPF(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/service-orders/code/100", nil)
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestServiceOrderHandler_FindByCode_CodigoInvalido(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/service-orders/code/abc?cpf=52998224725", nil)
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

// =============================================================================
// FindByCPF (rota publica)
// =============================================================================

func TestServiceOrderHandler_FindByCPF_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	orders := []*entities.ServiceOrder{entities.NewServiceOrder("c1", "v1")}
	m.listByCPF.EXPECT().Execute(gomock.Any(), "52998224725").Return(orders, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/service-orders/customer?cpf=52998224725", nil)
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestServiceOrderHandler_FindByCPF_SemCPF(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/service-orders/customer", nil)
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestServiceOrderHandler_FindByCPF_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	m.listByCPF.EXPECT().Execute(gomock.Any(), "00000000000").Return(nil, domainerrors.ErrNotFound)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/service-orders/customer?cpf=00000000000", nil)
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// =============================================================================
// Approve (rota publica)
// =============================================================================

func TestServiceOrderHandler_Approve_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	m.approve.EXPECT().Execute(gomock.Any(), 100, "52998224725").Return(nil)

	body, _ := json.Marshal(map[string]any{"customer_cpf": "52998224725"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/service-orders/code/100/approve", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d: %s", w.Code, w.Body.String())
	}
}

func TestServiceOrderHandler_Approve_BadRequest_CodigoInvalido(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)

	body, _ := json.Marshal(map[string]any{"customer_cpf": "52998224725"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/service-orders/code/abc/approve", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestServiceOrderHandler_Approve_BadRequest_SemBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/service-orders/code/100/approve", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestServiceOrderHandler_Approve_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	m.approve.EXPECT().Execute(gomock.Any(), 999, "52998224725").Return(domainerrors.ErrNotFound)

	body, _ := json.Marshal(map[string]any{"customer_cpf": "52998224725"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/service-orders/code/999/approve", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestServiceOrderHandler_Approve_InvalidStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	m.approve.EXPECT().Execute(gomock.Any(), 100, "52998224725").Return(domainerrors.ErrInvalidStatus)

	body, _ := json.Marshal(map[string]any{"customer_cpf": "52998224725"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/service-orders/code/100/approve", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}
}

// =============================================================================
// Reject (rota publica)
// =============================================================================

func TestServiceOrderHandler_Reject_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	m.reject.EXPECT().Execute(gomock.Any(), 100, "52998224725").Return(nil)

	body, _ := json.Marshal(map[string]any{"customer_cpf": "52998224725"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/service-orders/code/100/reject", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d: %s", w.Code, w.Body.String())
	}
}

func TestServiceOrderHandler_Reject_BadRequest_CodigoInvalido(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)

	body, _ := json.Marshal(map[string]any{"customer_cpf": "52998224725"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/service-orders/code/abc/reject", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestServiceOrderHandler_Reject_BadRequest_SemBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/service-orders/code/100/reject", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestServiceOrderHandler_Reject_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	m.reject.EXPECT().Execute(gomock.Any(), 999, "52998224725").Return(domainerrors.ErrNotFound)

	body, _ := json.Marshal(map[string]any{"customer_cpf": "52998224725"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/service-orders/code/999/reject", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestServiceOrderHandler_Reject_InvalidStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	m.reject.EXPECT().Execute(gomock.Any(), 100, "52998224725").Return(domainerrors.ErrInvalidStatus)

	body, _ := json.Marshal(map[string]any{"customer_cpf": "52998224725"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/service-orders/code/100/reject", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}
}
