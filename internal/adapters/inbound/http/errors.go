package httputil

import (
	"errors"
	"net/http"

	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func HandleError(c *gin.Context, err error, msg ...string) {
	code := statusFor(err)
	var message string
	if code == http.StatusInternalServerError {
		message = "internal server error"
	} else if len(msg) > 0 {
		message = msg[0]
	} else {
		message = err.Error()
	}
	c.JSON(code, ErrorResponse{
		Status:  http.StatusText(code),
		Message: message,
	})
}

func HandleBadRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Status:  http.StatusText(http.StatusBadRequest),
		Message: err.Error(),
	})
}

func statusFor(err error) int {
	switch {
	case errors.Is(err, domainerrors.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, domainerrors.ErrAlreadyExists):
		return http.StatusConflict
	case errors.Is(err, domainerrors.ErrInvalidCredentials),
		errors.Is(err, domainerrors.ErrInvalidRefreshToken):
		return http.StatusUnauthorized
	case errors.Is(err, domainerrors.ErrStatusNotAllowedForRequester),
		errors.Is(err, domainerrors.ErrDocumentMismatch):
		return http.StatusForbidden
	case errors.Is(err, domainerrors.ErrInvalidDocument),
		errors.Is(err, domainerrors.ErrInvalidPlate),
		errors.Is(err, domainerrors.ErrInvalidStatus),
		errors.Is(err, domainerrors.ErrInvalidStatusValue),
		errors.Is(err, domainerrors.ErrInsufficientStock),
		errors.Is(err, domainerrors.ErrOrderNotCancellable),
		errors.Is(err, domainerrors.ErrVehicleNotFromRequester):
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}
