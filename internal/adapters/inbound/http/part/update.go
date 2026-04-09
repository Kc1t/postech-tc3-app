package parthandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Update godoc
// @Summary     Atualizar peca/insumo
// @Tags        parts
// @Accept      json
// @Produce     json
// @Param       id path string true "Part ID"
// @Security    BearerAuth
// @Router      /parts/{id} [put]
func (h *PartHandler) Update(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
