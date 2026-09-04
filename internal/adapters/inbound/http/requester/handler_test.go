package requesterhandler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestRouter(h *RequesterHandler) *gin.Engine {
	r := gin.New()
	h.SetupRoutes(r.Group(""))
	return r
}

func newHandler(ctrl *gomock.Controller) (
	*RequesterHandler,
	*mocks.MockCreateRequesterUseCase,
	*mocks.MockGetRequesterUseCase,
	*mocks.MockGetRequesterByDocumentUseCase,
	*mocks.MockListRequestersUseCase,
	*mocks.MockUpdateRequesterUseCase,
	*mocks.MockDeleteRequesterUseCase,
	*mocks.MockListVehiclesByRequesterUseCase,
) {
	create := mocks.NewMockCreateRequesterUseCase(ctrl)
	getByID := mocks.NewMockGetRequesterUseCase(ctrl)
	getByDoc := mocks.NewMockGetRequesterByDocumentUseCase(ctrl)
	list := mocks.NewMockListRequestersUseCase(ctrl)
	update := mocks.NewMockUpdateRequesterUseCase(ctrl)
	del := mocks.NewMockDeleteRequesterUseCase(ctrl)
	listVehicles := mocks.NewMockListVehiclesByRequesterUseCase(ctrl)
	h := NewRequesterHandler(create, getByID, getByDoc, list, update, del, listVehicles)
	return h, create, getByID, getByDoc, list, update, del, listVehicles
}

// --- Create ---

func TestRequesterHandler_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, create, _, _, _, _, _, _ := newHandler(ctrl)
	create.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil)

	body, _ := json.Marshal(map[string]string{
		"name": "João Silva", "document": "52998224725",
		"email": "j@j.com", "phone": "11999",
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/requesters", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestRequesterHandler_Create_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, _, _, _, _, _ := newHandler(ctrl)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/requesters", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestRequesterHandler_Create_AlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, create, _, _, _, _, _, _ := newHandler(ctrl)
	create.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(domainerrors.ErrAlreadyExists)

	body, _ := json.Marshal(map[string]string{
		"name": "João", "document": "52998224725",
		"email": "j@j.com", "phone": "11999",
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/requesters", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", w.Code)
	}
}

// --- FindByID ---

func TestRequesterHandler_FindByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	requester := entities.ReconstituteRequester("", "João", "123", "j@j.com", "11999", time.Now(), time.Now())
	h, _, getByID, _, _, _, _, _ := newHandler(ctrl)
	getByID.EXPECT().Execute(gomock.Any(), "cust-1").Return(requester, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/requesters/cust-1", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestRequesterHandler_FindByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, getByID, _, _, _, _, _ := newHandler(ctrl)
	getByID.EXPECT().Execute(gomock.Any(), "cust-x").Return(nil, domainerrors.ErrNotFound)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/requesters/cust-x", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// --- FindByDocument ---

func TestRequesterHandler_FindByDocument_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	requester := entities.ReconstituteRequester("", "João", "12345678901", "j@j.com", "11999", time.Now(), time.Now())
	h, _, _, getByDoc, _, _, _, _ := newHandler(ctrl)
	getByDoc.EXPECT().Execute(gomock.Any(), "12345678901").Return(requester, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/requesters/document/12345678901", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestRequesterHandler_FindByDocument_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, getByDoc, _, _, _, _ := newHandler(ctrl)
	getByDoc.EXPECT().Execute(gomock.Any(), "00000000000").Return(nil, domainerrors.ErrNotFound)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/requesters/document/00000000000", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// --- FindAll ---

func TestRequesterHandler_FindAll_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	requesters := []*entities.Requester{entities.ReconstituteRequester("", "João", "123", "j@j.com", "11999", time.Now(), time.Now())}
	h, _, _, _, list, _, _, _ := newHandler(ctrl)
	list.EXPECT().Execute(gomock.Any()).Return(requesters, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/requesters", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestRequesterHandler_FindAll_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, _, list, _, _, _ := newHandler(ctrl)
	list.EXPECT().Execute(gomock.Any()).Return(nil, errors.New("db error"))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/requesters", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// --- Update ---

func TestRequesterHandler_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	requester := entities.ReconstituteRequester("", "João", "123", "j@j.com", "11999", time.Now(), time.Now())
	h, _, getByID, _, _, update, _, _ := newHandler(ctrl)
	getByID.EXPECT().Execute(gomock.Any(), "cust-1").Return(requester, nil)
	update.EXPECT().Execute(gomock.Any(), requester).Return(nil)

	body, _ := json.Marshal(map[string]string{"name": "João Novo"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/requesters/cust-1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestRequesterHandler_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, getByID, _, _, _, _, _ := newHandler(ctrl)
	getByID.EXPECT().Execute(gomock.Any(), "cust-x").Return(nil, domainerrors.ErrNotFound)

	body, _ := json.Marshal(map[string]string{"name": "Novo"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/requesters/cust-x", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestRequesterHandler_Update_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, _, _, _, _, _ := newHandler(ctrl)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/requesters/cust-1", bytes.NewBufferString(`invalid`))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestRequesterHandler_Update_UpdateFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	requester := entities.ReconstituteRequester("", "João", "123", "j@j.com", "11999", time.Now(), time.Now())
	h, _, getByID, _, _, update, _, _ := newHandler(ctrl)
	getByID.EXPECT().Execute(gomock.Any(), "cust-1").Return(requester, nil)
	update.EXPECT().Execute(gomock.Any(), requester).Return(errors.New("db error"))

	body, _ := json.Marshal(map[string]string{"name": "João Novo"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/requesters/cust-1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// --- Delete ---

func TestRequesterHandler_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, _, _, _, del, _ := newHandler(ctrl)
	del.EXPECT().Execute(gomock.Any(), "cust-1").Return(nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/requesters/cust-1", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func TestRequesterHandler_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, _, _, _, del, _ := newHandler(ctrl)
	del.EXPECT().Execute(gomock.Any(), "cust-x").Return(domainerrors.ErrNotFound)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/requesters/cust-x", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// --- ListVehicles (GET /requesters/:id/vehicles) ---

func TestRequesterHandler_ListVehicles_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, _, _, _, _, listVehicles := newHandler(ctrl)
	now := time.Now()
	v1 := entities.ReconstituteVehicle("v1", "cust-1", "ABC1234", "Fiat", "Strada", 2021, now, now)
	v2 := entities.ReconstituteVehicle("v2", "cust-1", "XYZ9A88", "VW", "Gol", 2019, now, now)
	listVehicles.EXPECT().Execute(gomock.Any(), "cust-1").Return([]*entities.Vehicle{v1, v2}, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/requesters/cust-1/vehicles", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var body []map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("json decode: %v", err)
	}
	if len(body) != 2 {
		t.Errorf("retornou %d veiculos, esperava 2", len(body))
	}
}

func TestRequesterHandler_ListVehicles_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, _, _, _, _, listVehicles := newHandler(ctrl)
	listVehicles.EXPECT().Execute(gomock.Any(), "cust-sem-veh").Return([]*entities.Vehicle{}, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/requesters/cust-sem-veh/vehicles", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if w.Body.String() != "[]" {
		t.Errorf("body = %q, esperava \"[]\"", w.Body.String())
	}
}

func TestRequesterHandler_ListVehicles_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, _, _, _, _, listVehicles := newHandler(ctrl)
	listVehicles.EXPECT().Execute(gomock.Any(), "cust-1").Return(nil, errors.New("db unavailable"))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/requesters/cust-1/vehicles", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
