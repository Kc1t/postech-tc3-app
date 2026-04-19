package customerhandler

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

func newTestRouter(h *CustomerHandler) *gin.Engine {
	r := gin.New()
	h.SetupRoutes(r.Group(""))
	return r
}

func newHandler(ctrl *gomock.Controller) (
	*CustomerHandler,
	*mocks.MockCreateCustomerUseCase,
	*mocks.MockGetCustomerUseCase,
	*mocks.MockGetCustomerByDocumentUseCase,
	*mocks.MockListCustomersUseCase,
	*mocks.MockUpdateCustomerUseCase,
	*mocks.MockDeleteCustomerUseCase,
) {
	create := mocks.NewMockCreateCustomerUseCase(ctrl)
	getByID := mocks.NewMockGetCustomerUseCase(ctrl)
	getByDoc := mocks.NewMockGetCustomerByDocumentUseCase(ctrl)
	list := mocks.NewMockListCustomersUseCase(ctrl)
	update := mocks.NewMockUpdateCustomerUseCase(ctrl)
	del := mocks.NewMockDeleteCustomerUseCase(ctrl)
	h := NewCustomerHandler(create, getByID, getByDoc, list, update, del)
	return h, create, getByID, getByDoc, list, update, del
}

// --- Create ---

func TestCustomerHandler_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, create, _, _, _, _, _ := newHandler(ctrl)
	create.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil)

	body, _ := json.Marshal(map[string]string{
		"name": "João Silva", "document": "52998224725",
		"email": "j@j.com", "phone": "11999",
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCustomerHandler_Create_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, _, _, _, _ := newHandler(ctrl)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCustomerHandler_Create_AlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, create, _, _, _, _, _ := newHandler(ctrl)
	create.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(domainerrors.ErrAlreadyExists)

	body, _ := json.Marshal(map[string]string{
		"name": "João", "document": "52998224725",
		"email": "j@j.com", "phone": "11999",
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", w.Code)
	}
}

// --- FindByID ---

func TestCustomerHandler_FindByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customer := entities.ReconstituteCustomer("", "João", "123", "j@j.com", "11999", time.Now(), time.Now())
	h, _, getByID, _, _, _, _ := newHandler(ctrl)
	getByID.EXPECT().Execute(gomock.Any(), "cust-1").Return(customer, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/customers/cust-1", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestCustomerHandler_FindByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, getByID, _, _, _, _ := newHandler(ctrl)
	getByID.EXPECT().Execute(gomock.Any(), "cust-x").Return(nil, domainerrors.ErrNotFound)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/customers/cust-x", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// --- FindByDocument ---

func TestCustomerHandler_FindByDocument_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customer := entities.ReconstituteCustomer("", "João", "12345678901", "j@j.com", "11999", time.Now(), time.Now())
	h, _, _, getByDoc, _, _, _ := newHandler(ctrl)
	getByDoc.EXPECT().Execute(gomock.Any(), "12345678901").Return(customer, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/customers/document/12345678901", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestCustomerHandler_FindByDocument_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, getByDoc, _, _, _ := newHandler(ctrl)
	getByDoc.EXPECT().Execute(gomock.Any(), "00000000000").Return(nil, domainerrors.ErrNotFound)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/customers/document/00000000000", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// --- FindAll ---

func TestCustomerHandler_FindAll_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customers := []*entities.Customer{entities.ReconstituteCustomer("", "João", "123", "j@j.com", "11999", time.Now(), time.Now())}
	h, _, _, _, list, _, _ := newHandler(ctrl)
	list.EXPECT().Execute(gomock.Any()).Return(customers, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/customers", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestCustomerHandler_FindAll_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, _, list, _, _ := newHandler(ctrl)
	list.EXPECT().Execute(gomock.Any()).Return(nil, errors.New("db error"))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/customers", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// --- Update ---

func TestCustomerHandler_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customer := entities.ReconstituteCustomer("", "João", "123", "j@j.com", "11999", time.Now(), time.Now())
	h, _, getByID, _, _, update, _ := newHandler(ctrl)
	getByID.EXPECT().Execute(gomock.Any(), "cust-1").Return(customer, nil)
	update.EXPECT().Execute(gomock.Any(), customer).Return(nil)

	body, _ := json.Marshal(map[string]string{"name": "João Novo"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/customers/cust-1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCustomerHandler_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, getByID, _, _, _, _ := newHandler(ctrl)
	getByID.EXPECT().Execute(gomock.Any(), "cust-x").Return(nil, domainerrors.ErrNotFound)

	body, _ := json.Marshal(map[string]string{"name": "Novo"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/customers/cust-x", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestCustomerHandler_Update_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, _, _, _, _ := newHandler(ctrl)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/customers/cust-1", bytes.NewBufferString(`invalid`))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCustomerHandler_Update_UpdateFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customer := entities.ReconstituteCustomer("", "João", "123", "j@j.com", "11999", time.Now(), time.Now())
	h, _, getByID, _, _, update, _ := newHandler(ctrl)
	getByID.EXPECT().Execute(gomock.Any(), "cust-1").Return(customer, nil)
	update.EXPECT().Execute(gomock.Any(), customer).Return(errors.New("db error"))

	body, _ := json.Marshal(map[string]string{"name": "João Novo"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/customers/cust-1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// --- Delete ---

func TestCustomerHandler_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, _, _, _, del := newHandler(ctrl)
	del.EXPECT().Execute(gomock.Any(), "cust-1").Return(nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/customers/cust-1", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func TestCustomerHandler_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h, _, _, _, _, _, del := newHandler(ctrl)
	del.EXPECT().Execute(gomock.Any(), "cust-x").Return(domainerrors.ErrNotFound)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/customers/cust-x", nil)
	newTestRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
