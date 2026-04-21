package httputil

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newGinContext(w *httptest.ResponseRecorder) *gin.Context {
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	return c
}

func TestHandleError_NotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c := newGinContext(w)
	HandleError(c, domainerrors.ErrNotFound)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Status != http.StatusText(http.StatusNotFound) {
		t.Errorf("expected status %q, got %q", http.StatusText(http.StatusNotFound), resp.Status)
	}
}

func TestHandleError_AlreadyExists(t *testing.T) {
	w := httptest.NewRecorder()
	c := newGinContext(w)
	HandleError(c, domainerrors.ErrAlreadyExists)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", w.Code)
	}
}

func TestHandleError_InvalidCredentials(t *testing.T) {
	w := httptest.NewRecorder()
	c := newGinContext(w)
	HandleError(c, domainerrors.ErrInvalidCredentials)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestHandleError_InvalidRefreshToken(t *testing.T) {
	w := httptest.NewRecorder()
	c := newGinContext(w)
	HandleError(c, domainerrors.ErrInvalidRefreshToken)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestHandleError_StatusNotAllowedForCustomer(t *testing.T) {
	w := httptest.NewRecorder()
	c := newGinContext(w)
	HandleError(c, domainerrors.ErrStatusNotAllowedForCustomer)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestHandleError_UnprocessableEntity_Errors(t *testing.T) {
	unprocessable := []error{
		domainerrors.ErrInvalidDocument,
		domainerrors.ErrInvalidPlate,
		domainerrors.ErrInvalidStatus,
		domainerrors.ErrInsufficientStock,
		domainerrors.ErrOrderNotCancellable,
	}
	for _, e := range unprocessable {
		w := httptest.NewRecorder()
		c := newGinContext(w)
		HandleError(c, e)
		if w.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected 422 for %v, got %d", e, w.Code)
		}
	}
}

func TestHandleError_InternalServerError(t *testing.T) {
	w := httptest.NewRecorder()
	c := newGinContext(w)
	HandleError(c, errors.New("unexpected db error"))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Message != "internal server error" {
		t.Errorf("expected masked message, got %q", resp.Message)
	}
}

func TestHandleError_CustomMessage_4xx(t *testing.T) {
	w := httptest.NewRecorder()
	c := newGinContext(w)
	HandleError(c, domainerrors.ErrNotFound, "custom not found msg")

	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Message != "custom not found msg" {
		t.Errorf("expected custom message, got %q", resp.Message)
	}
}

func TestHandleError_CustomMessage_5xx_Masked(t *testing.T) {
	w := httptest.NewRecorder()
	c := newGinContext(w)
	HandleError(c, errors.New("secret db error"), "this should be ignored")

	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Message != "internal server error" {
		t.Errorf("expected masked 5xx message, got %q", resp.Message)
	}
}

func TestHandleBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	c := newGinContext(w)
	HandleBadRequest(c, errors.New("validation failed"))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Status != http.StatusText(http.StatusBadRequest) {
		t.Errorf("expected status %q, got %q", http.StatusText(http.StatusBadRequest), resp.Status)
	}
	if resp.Message != "validation failed" {
		t.Errorf("expected message %q, got %q", "validation failed", resp.Message)
	}
}
