package servicehandler

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

func newTestRouter(h *ServiceHandler) *gin.Engine {
	r := gin.New()
	h.SetupRoutes(r.Group(""))
	return r
}

func newHandler(ctrl *gomock.Controller) (
	*ServiceHandler,
	*mocks.MockCreateServiceUseCase,
	*mocks.MockGetServiceUseCase,
	*mocks.MockListServicesUseCase,
	*mocks.MockUpdateServiceUseCase,
	*mocks.MockDeleteServiceUseCase,
) {
	create := mocks.NewMockCreateServiceUseCase(ctrl)
	getByID := mocks.NewMockGetServiceUseCase(ctrl)
	list := mocks.NewMockListServicesUseCase(ctrl)
	update := mocks.NewMockUpdateServiceUseCase(ctrl)
	del := mocks.NewMockDeleteServiceUseCase(ctrl)
	h := NewServiceHandler(create, getByID, list, update, del)
	return h, create, getByID, list, update, del
}

// --- Create ---

func TestServiceHandler_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, create, _, _, _, _ := newHandler(ctrl)
	create.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil)

	body, _ := json.Marshal(map[string]any{
		"name": "Troca de óleo", "price": 150.0, "duration_min": 60,
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/services", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestServiceHandler_Create_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, _, _, _ := newHandler(ctrl)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/services", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestServiceHandler_Create_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, create, _, _, _, _ := newHandler(ctrl)
	create.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(errors.New("db error"))

	body, _ := json.Marshal(map[string]any{
		"name": "Troca de óleo", "price": 150.0, "duration_min": 60,
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/services", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// --- FindByID ---

func TestServiceHandler_FindByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := entities.NewService("Troca de óleo", "Desc", 150.0, 60)
	h, _, getByID, _, _, _ := newHandler(ctrl)
	getByID.EXPECT().Execute(gomock.Any(), "svc-1").Return(svc, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/services/svc-1", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestServiceHandler_FindByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, getByID, _, _, _ := newHandler(ctrl)
	getByID.EXPECT().Execute(gomock.Any(), "svc-x").Return(nil, domainerrors.ErrNotFound)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/services/svc-x", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// --- FindAll ---

func TestServiceHandler_FindAll_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svcs := []*entities.Service{entities.NewService("Troca de óleo", "Desc", 150.0, 60)}
	h, _, _, list, _, _ := newHandler(ctrl)
	list.EXPECT().Execute(gomock.Any()).Return(svcs, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/services", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestServiceHandler_FindAll_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, list, _, _ := newHandler(ctrl)
	list.EXPECT().Execute(gomock.Any()).Return(nil, errors.New("db error"))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/services", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// --- Update ---

func TestServiceHandler_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := entities.NewService("Troca de óleo", "Desc", 150.0, 60)
	h, _, getByID, _, update, _ := newHandler(ctrl)
	getByID.EXPECT().Execute(gomock.Any(), "svc-1").Return(svc, nil)
	update.EXPECT().Execute(gomock.Any(), svc).Return(nil)

	body, _ := json.Marshal(map[string]any{"name": "Alinhamento", "price": 200.0})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/services/svc-1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestServiceHandler_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, getByID, _, _, _ := newHandler(ctrl)
	getByID.EXPECT().Execute(gomock.Any(), "svc-x").Return(nil, domainerrors.ErrNotFound)

	body, _ := json.Marshal(map[string]any{"name": "X"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/services/svc-x", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestServiceHandler_Update_UpdateFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := entities.NewService("Troca de óleo", "Desc", 150.0, 60)
	h, _, getByID, _, update, _ := newHandler(ctrl)
	getByID.EXPECT().Execute(gomock.Any(), "svc-1").Return(svc, nil)
	update.EXPECT().Execute(gomock.Any(), svc).Return(errors.New("db error"))

	body, _ := json.Marshal(map[string]any{"name": "Alinhamento"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/services/svc-1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// --- Delete ---

func TestServiceHandler_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, _, _, del := newHandler(ctrl)
	del.EXPECT().Execute(gomock.Any(), "svc-1").Return(nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/services/svc-1", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func TestServiceHandler_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, _, _, del := newHandler(ctrl)
	del.EXPECT().Execute(gomock.Any(), "svc-x").Return(domainerrors.ErrNotFound)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/services/svc-x", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
