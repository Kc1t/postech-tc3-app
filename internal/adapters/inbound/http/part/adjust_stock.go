package parthandler

import (
	"errors"
	"net/http"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/gin-gonic/gin"
)

// AdjustStock godoc
// @Summary     Ajustar estoque de peca/insumo
// @Tags        parts
// @Accept      json
// @Produce     json
// @Param       id path string true "Part ID"
// @Param       body body commands.AdjustStockRequest true "Delta de estoque"
// @Success     204
// @Security    BearerAuth
// @Router      /parts/{id}/stock [patch]
func (h *PartHandler) AdjustStock(c *gin.Context) {
	var req commands.AdjustStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.HandleBadRequest(c, err)
		return
	}

	if err := h.adjustStock.Execute(c.Request.Context(), c.Param("id"), req.Delta); err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrNotFound):
			httputil.HandleError(c, err, msgNotFound)
		case errors.Is(err, domainerrors.ErrInsufficientStock):
			httputil.HandleError(c, err, msgInsufficientStock)
		default:
			httputil.HandleError(c, err)
		}
		return
	}

	c.Status(http.StatusNoContent)
}
