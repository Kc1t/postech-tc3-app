package parthandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdjustStock godoc
// @Summary     Ajustar estoque de peca/insumo
// @Tags        parts
// @Accept      json
// @Produce     json
// @Param       id path string true "Part ID"
// @Security    BearerAuth
// @Router      /parts/{id}/stock [patch]
func (h *PartHandler) AdjustStock(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
