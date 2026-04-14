package vehiclehandler

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

func newTestRouter(h *VehicleHandler) *gin.Engine {
	r := gin.New()
	h.SetupRoutes(r.Group(""))
	return r
}

func newHandler(ctrl *gomock.Controller) (
	*VehicleHandler,
	*mocks.MockCreateVehicleUseCase,
	*mocks.MockGetVehicleUseCase,
	*mocks.MockListVehiclesUseCase,
	*mocks.MockUpdateVehicleUseCase,
	*mocks.MockDeleteVehicleUseCase,
) {
	create := mocks.NewMockCreateVehicleUseCase(ctrl)
	getByID := mocks.NewMockGetVehicleUseCase(ctrl)
	list := mocks.NewMockListVehiclesUseCase(ctrl)
	update := mocks.NewMockUpdateVehicleUseCase(ctrl)
	del := mocks.NewMockDeleteVehicleUseCase(ctrl)
	h := NewVehicleHandler(create, getByID, list, update, del)
	return h, create, getByID, list, update, del
}

// --- Create ---

func TestVehicleHandler_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, create, _, _, _, _ := newHandler(ctrl)
	create.EXPECT().Execute(gomock.Any(), "12345678901", gomock.Any()).Return(nil)

	body, _ := json.Marshal(map[string]any{
		"customer_document": "12345678901",
		"plate": "ABC1234", "brand": "Toyota", "model": "Corolla", "year": 2020,
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/vehicles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestVehicleHandler_Create_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, _, _, _ := newHandler(ctrl)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/vehicles", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestVehicleHandler_Create_CustomerNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, create, _, _, _, _ := newHandler(ctrl)
	create.EXPECT().Execute(gomock.Any(), gomock.Any(), gomock.Any()).Return(domainerrors.ErrNotFound)

	body, _ := json.Marshal(map[string]any{
		"customer_document": "00000000000",
		"plate": "ABC1234", "brand": "Toyota", "model": "Corolla", "year": 2020,
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/vehicles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestVehicleHandler_Create_AlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, create, _, _, _, _ := newHandler(ctrl)
	create.EXPECT().Execute(gomock.Any(), gomock.Any(), gomock.Any()).Return(domainerrors.ErrAlreadyExists)

	body, _ := json.Marshal(map[string]any{
		"customer_document": "12345678901",
		"plate": "ABC1234", "brand": "Toyota", "model": "Corolla", "year": 2020,
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/vehicles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", w.Code)
	}
}

func TestVehicleHandler_Create_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, create, _, _, _, _ := newHandler(ctrl)
	create.EXPECT().Execute(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("unexpected"))

	body, _ := json.Marshal(map[string]any{
		"customer_document": "12345678901",
		"plate": "ABC1234", "brand": "Toyota", "model": "Corolla", "year": 2020,
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/vehicles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// --- FindByID ---

func TestVehicleHandler_FindByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vehicle := entities.NewVehicle("cust-1", "ABC1234", "Toyota", "Corolla", 2020)
	h, _, getByID, _, _, _ := newHandler(ctrl)
	getByID.EXPECT().Execute(gomock.Any(), "veh-1").Return(vehicle, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/vehicles/veh-1", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestVehicleHandler_FindByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, getByID, _, _, _ := newHandler(ctrl)
	getByID.EXPECT().Execute(gomock.Any(), "veh-x").Return(nil, domainerrors.ErrNotFound)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/vehicles/veh-x", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// --- FindAll ---

func TestVehicleHandler_FindAll_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vehicles := []*entities.Vehicle{entities.NewVehicle("cust-1", "ABC1234", "Toyota", "Corolla", 2020)}
	h, _, _, list, _, _ := newHandler(ctrl)
	list.EXPECT().Execute(gomock.Any()).Return(vehicles, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/vehicles", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestVehicleHandler_FindAll_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, list, _, _ := newHandler(ctrl)
	list.EXPECT().Execute(gomock.Any()).Return(nil, errors.New("db error"))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/vehicles", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// --- Update ---

func TestVehicleHandler_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vehicle := entities.NewVehicle("cust-1", "ABC1234", "Toyota", "Corolla", 2020)
	h, _, getByID, _, update, _ := newHandler(ctrl)
	getByID.EXPECT().Execute(gomock.Any(), "veh-1").Return(vehicle, nil)
	update.EXPECT().Execute(gomock.Any(), vehicle).Return(nil)

	body, _ := json.Marshal(map[string]any{"brand": "Honda"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/vehicles/veh-1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestVehicleHandler_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, getByID, _, _, _ := newHandler(ctrl)
	getByID.EXPECT().Execute(gomock.Any(), "veh-x").Return(nil, domainerrors.ErrNotFound)

	body, _ := json.Marshal(map[string]any{"brand": "Honda"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/vehicles/veh-x", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestVehicleHandler_Update_UpdateFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vehicle := entities.NewVehicle("cust-1", "ABC1234", "Toyota", "Corolla", 2020)
	h, _, getByID, _, update, _ := newHandler(ctrl)
	getByID.EXPECT().Execute(gomock.Any(), "veh-1").Return(vehicle, nil)
	update.EXPECT().Execute(gomock.Any(), vehicle).Return(errors.New("db error"))

	body, _ := json.Marshal(map[string]any{"brand": "Honda"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/vehicles/veh-1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// --- Delete ---

func TestVehicleHandler_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, _, _, del := newHandler(ctrl)
	del.EXPECT().Execute(gomock.Any(), "veh-1").Return(nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/vehicles/veh-1", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func TestVehicleHandler_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, _, _, del := newHandler(ctrl)
	del.EXPECT().Execute(gomock.Any(), "veh-x").Return(domainerrors.ErrNotFound)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/vehicles/veh-x", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
