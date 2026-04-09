package parthandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// FindAll godoc
// @Summary     Listar pecas/insumos
// @Tags        parts
// @Produce     json
// @Security    BearerAuth
// @Router      /parts [get]
func (h *PartHandler) FindAll(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
