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
	handler              *ServiceOrderHandler
	create               *mocks.MockCreateServiceOrderUseCase
	getByID              *mocks.MockGetServiceOrderUseCase
	getByCode            *mocks.MockGetServiceOrderByCodeUseCase
	listAll              *mocks.MockListServiceOrdersUseCase
	listByDocument       *mocks.MockListServiceOrdersByDocumentUseCase
	updateStatus         *mocks.MockUpdateServiceOrderStatusUseCase
	update               *mocks.MockUpdateServiceOrderUseCase
	del                  *mocks.MockDeleteServiceOrderUseCase
	updateStatusByCode   *mocks.MockUpdateServiceOrderStatusByCodeUseCase
	averageExecutionTime *mocks.MockGetAverageExecutionTimeUseCase
}

func newHandler(ctrl *gomock.Controller) handlerMocks {
	m := handlerMocks{
		create:               mocks.NewMockCreateServiceOrderUseCase(ctrl),
		getByID:              mocks.NewMockGetServiceOrderUseCase(ctrl),
		getByCode:            mocks.NewMockGetServiceOrderByCodeUseCase(ctrl),
		listAll:              mocks.NewMockListServiceOrdersUseCase(ctrl),
		listByDocument:       mocks.NewMockListServiceOrdersByDocumentUseCase(ctrl),
		updateStatus:         mocks.NewMockUpdateServiceOrderStatusUseCase(ctrl),
		update:               mocks.NewMockUpdateServiceOrderUseCase(ctrl),
		del:                  mocks.NewMockDeleteServiceOrderUseCase(ctrl),
		updateStatusByCode:   mocks.NewMockUpdateServiceOrderStatusByCodeUseCase(ctrl),
		averageExecutionTime: mocks.NewMockGetAverageExecutionTimeUseCase(ctrl),
	}
	m.handler = NewServiceOrderHandler(
		m.create, m.getByID, m.getByCode, m.listAll, m.listByDocument,
		m.updateStatus, m.update, m.del, m.updateStatusByCode,
		m.averageExecutionTime,
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
		"customer_document":  "52998224725",
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
		"customer_document":  "52998224725",
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
		"customer_document":  "52998224725",
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
		"customer_document":  "52998224725",
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
	req := httptest.NewRequest(http.MethodGet, "/service-orders/code/100?document=52998224725", nil)
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
	req := httptest.NewRequest(http.MethodGet, "/service-orders/code/999?document=52998224725", nil)
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestServiceOrderHandler_FindByCode_SemDocument(t *testing.T) {
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
	req := httptest.NewRequest(http.MethodGet, "/service-orders/code/abc?document=52998224725", nil)
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

// =============================================================================
// FindByDocument (rota publica)
// =============================================================================

func TestServiceOrderHandler_FindByDocument_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	orders := []*entities.ServiceOrder{entities.NewServiceOrder("c1", "v1")}
	m.listByDocument.EXPECT().Execute(gomock.Any(), "52998224725").Return(orders, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/service-orders/customer?document=52998224725", nil)
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestServiceOrderHandler_FindByDocument_SemDocument(t *testing.T) {
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

func TestServiceOrderHandler_FindByDocument_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	m.listByDocument.EXPECT().Execute(gomock.Any(), "52998224725").Return(nil, domainerrors.ErrNotFound)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/service-orders/customer?document=52998224725", nil)
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// =============================================================================
// UpdateStatusByCode (rota publica)
// =============================================================================

func TestServiceOrderHandler_UpdateStatusByCode_Aprovacao(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	m.updateStatusByCode.EXPECT().Execute(gomock.Any(), 100, "52998224725", entities.StatusInExecution).Return(nil)

	body, _ := json.Marshal(map[string]any{"customer_document": "52998224725", "status": "in_execution"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/service-orders/code/100/status", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d: %s", w.Code, w.Body.String())
	}
}

func TestServiceOrderHandler_UpdateStatusByCode_Rejeicao(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	m.updateStatusByCode.EXPECT().Execute(gomock.Any(), 100, "52998224725", entities.StatusReceived).Return(nil)

	body, _ := json.Marshal(map[string]any{"customer_document": "52998224725", "status": "received"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/service-orders/code/100/status", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d: %s", w.Code, w.Body.String())
	}
}

func TestServiceOrderHandler_UpdateStatusByCode_CodigoInvalido(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)

	body, _ := json.Marshal(map[string]any{"customer_document": "52998224725", "status": "in_execution"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/service-orders/code/abc/status", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestServiceOrderHandler_UpdateStatusByCode_SemBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/service-orders/code/100/status", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestServiceOrderHandler_UpdateStatusByCode_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	m.updateStatusByCode.EXPECT().Execute(gomock.Any(), 999, "52998224725", entities.StatusInExecution).Return(domainerrors.ErrNotFound)

	body, _ := json.Marshal(map[string]any{"customer_document": "52998224725", "status": "in_execution"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/service-orders/code/999/status", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestServiceOrderHandler_UpdateStatusByCode_StatusNaoPermitido(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	m.updateStatusByCode.EXPECT().
		Execute(gomock.Any(), 100, "52998224725", entities.StatusFinished).
		Return(domainerrors.ErrStatusNotAllowedForCustomer)

	body, _ := json.Marshal(map[string]any{"customer_document": "52998224725", "status": "finished"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/service-orders/code/100/status", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d: %s", w.Code, w.Body.String())
	}
}

func TestServiceOrderHandler_UpdateStatusByCode_TransicaoInvalida(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newHandler(ctrl)
	m.updateStatusByCode.EXPECT().Execute(gomock.Any(), 100, "52998224725", entities.StatusInExecution).Return(domainerrors.ErrInvalidStatus)

	body, _ := json.Marshal(map[string]any{"customer_document": "52998224725", "status": "in_execution"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/service-orders/code/100/status", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(m.handler).ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}
}
