package parthandler

import (
	"errors"
	"net/http"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.adjustStock.Execute(c.Request.Context(), c.Param("id"), req.Delta); err != nil {
		if errors.Is(err, domainerrors.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "part not found"})
			return
		}
		if errors.Is(err, domainerrors.ErrInsufficientStock) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "insufficient stock"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
