package httputil

import (
	"errors"
	"net/http"

	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/gin-gonic/gin"
)

// ErrorResponse é a struct padrão de resposta de erro da API.
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// HandleError mapeia erros de domínio para o HTTP status correto e escreve a resposta.
// msg é opcional: quando fornecida, substitui err.Error() apenas em respostas 4xx.
// Erros 5xx sempre retornam "internal server error" — detalhes de infra nunca são expostos.
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

// HandleBadRequest escreve uma resposta 400 padronizada para erros de bind/validação.
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
	case errors.Is(err, domainerrors.ErrStatusNotAllowedForCustomer):
		return http.StatusForbidden
	case errors.Is(err, domainerrors.ErrInvalidDocument),
		errors.Is(err, domainerrors.ErrInvalidPlate),
		errors.Is(err, domainerrors.ErrInvalidStatus),
		errors.Is(err, domainerrors.ErrInvalidStatusValue),
		errors.Is(err, domainerrors.ErrInsufficientStock),
		errors.Is(err, domainerrors.ErrOrderNotCancellable),
		errors.Is(err, domainerrors.ErrVehicleNotFromCustomer):
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}
