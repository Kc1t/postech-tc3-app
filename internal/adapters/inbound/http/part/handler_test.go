package parthandler

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

func newTestRouter(h *PartHandler) *gin.Engine {
	r := gin.New()
	h.SetupRoutes(r.Group(""))
	return r
}

func newHandler(ctrl *gomock.Controller) (
	*PartHandler,
	*mocks.MockCreatePartUseCase,
	*mocks.MockGetPartUseCase,
	*mocks.MockListPartsUseCase,
	*mocks.MockUpdatePartUseCase,
	*mocks.MockDeletePartUseCase,
	*mocks.MockAdjustPartStockUseCase,
) {
	create := mocks.NewMockCreatePartUseCase(ctrl)
	getByID := mocks.NewMockGetPartUseCase(ctrl)
	list := mocks.NewMockListPartsUseCase(ctrl)
	update := mocks.NewMockUpdatePartUseCase(ctrl)
	del := mocks.NewMockDeletePartUseCase(ctrl)
	adjust := mocks.NewMockAdjustPartStockUseCase(ctrl)
	h := NewPartHandler(create, getByID, list, update, del, adjust)
	return h, create, getByID, list, update, del, adjust
}

// --- Create ---

func TestPartHandler_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, create, _, _, _, _, _ := newHandler(ctrl)
	create.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil)

	body, _ := json.Marshal(map[string]any{
		"name": "Filtro", "unit": "unidade", "price": 49.90, "stock": 10,
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/parts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestPartHandler_Create_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, _, _, _, _ := newHandler(ctrl)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/parts", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestPartHandler_Create_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, create, _, _, _, _, _ := newHandler(ctrl)
	create.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(errors.New("db error"))

	body, _ := json.Marshal(map[string]any{
		"name": "Filtro", "unit": "unidade", "price": 49.90, "stock": 10,
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/parts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// --- FindByID ---

func TestPartHandler_FindByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	part := entities.NewPart("Filtro", "Desc", "unidade", 49.90, 10)
	h, _, getByID, _, _, _, _ := newHandler(ctrl)
	getByID.EXPECT().Execute(gomock.Any(), "part-1").Return(part, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/parts/part-1", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestPartHandler_FindByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, getByID, _, _, _, _ := newHandler(ctrl)
	getByID.EXPECT().Execute(gomock.Any(), "part-x").Return(nil, domainerrors.ErrNotFound)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/parts/part-x", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// --- FindAll ---

func TestPartHandler_FindAll_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	parts := []*entities.Part{entities.NewPart("Filtro", "Desc", "unidade", 49.90, 10)}
	h, _, _, list, _, _, _ := newHandler(ctrl)
	list.EXPECT().Execute(gomock.Any()).Return(parts, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/parts", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestPartHandler_FindAll_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, list, _, _, _ := newHandler(ctrl)
	list.EXPECT().Execute(gomock.Any()).Return(nil, errors.New("db error"))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/parts", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// --- Update ---

func TestPartHandler_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	part := entities.NewPart("Filtro", "Desc", "unidade", 49.90, 10)
	h, _, getByID, _, update, _, _ := newHandler(ctrl)
	getByID.EXPECT().Execute(gomock.Any(), "part-1").Return(part, nil)
	update.EXPECT().Execute(gomock.Any(), part).Return(nil)

	body, _ := json.Marshal(map[string]any{"name": "Novo Filtro", "price": 59.90})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/parts/part-1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestPartHandler_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, getByID, _, _, _, _ := newHandler(ctrl)
	getByID.EXPECT().Execute(gomock.Any(), "part-x").Return(nil, domainerrors.ErrNotFound)

	body, _ := json.Marshal(map[string]any{"name": "X"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/parts/part-x", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestPartHandler_Update_UpdateFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	part := entities.NewPart("Filtro", "Desc", "unidade", 49.90, 10)
	h, _, getByID, _, update, _, _ := newHandler(ctrl)
	getByID.EXPECT().Execute(gomock.Any(), "part-1").Return(part, nil)
	update.EXPECT().Execute(gomock.Any(), part).Return(errors.New("db error"))

	body, _ := json.Marshal(map[string]any{"name": "Novo Filtro"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/parts/part-1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// --- Delete ---

func TestPartHandler_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, _, _, del, _ := newHandler(ctrl)
	del.EXPECT().Execute(gomock.Any(), "part-1").Return(nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/parts/part-1", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func TestPartHandler_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, _, _, del, _ := newHandler(ctrl)
	del.EXPECT().Execute(gomock.Any(), "part-x").Return(domainerrors.ErrNotFound)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/parts/part-x", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// --- AdjustStock ---

func TestPartHandler_AdjustStock_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, _, _, _, adjust := newHandler(ctrl)
	adjust.EXPECT().Execute(gomock.Any(), "part-1", 5).Return(nil)

	body, _ := json.Marshal(map[string]any{"delta": 5})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/parts/part-1/stock", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d: %s", w.Code, w.Body.String())
	}
}

func TestPartHandler_AdjustStock_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, _, _, _, _ := newHandler(ctrl)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/parts/part-1/stock", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestPartHandler_AdjustStock_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, _, _, _, adjust := newHandler(ctrl)
	adjust.EXPECT().Execute(gomock.Any(), "part-x", 5).Return(domainerrors.ErrNotFound)

	body, _ := json.Marshal(map[string]any{"delta": 5})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/parts/part-x/stock", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestPartHandler_AdjustStock_InsufficientStock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, _, _, _, adjust := newHandler(ctrl)
	adjust.EXPECT().Execute(gomock.Any(), "part-1", -100).Return(domainerrors.ErrInsufficientStock)

	body, _ := json.Marshal(map[string]any{"delta": -100})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/parts/part-1/stock", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}
}

func TestPartHandler_AdjustStock_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, _, _, _, adjust := newHandler(ctrl)
	adjust.EXPECT().Execute(gomock.Any(), "part-1", 5).Return(errors.New("unexpected"))

	body, _ := json.Marshal(map[string]any{"delta": 5})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/parts/part-1/stock", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
