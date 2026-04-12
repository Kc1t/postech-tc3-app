package parthandler

import (
	"net/http"

	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	"github.com/gin-gonic/gin"
)

// FindAll godoc
// @Summary     Listar pecas/insumos
// @Tags        parts
// @Produce     json
// @Success     200 {array} commands.PartResponse
// @Security    BearerAuth
// @Router      /parts [get]
func (h *PartHandler) FindAll(c *gin.Context) {
	parts, err := h.listAll.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, commands.ToPartListResponse(parts))
}
