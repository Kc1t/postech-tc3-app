package parthandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// FindByID godoc
// @Summary     Buscar peca/insumo por ID
// @Tags        parts
// @Produce     json
// @Param       id path string true "Part ID"
// @Security    BearerAuth
// @Router      /parts/{id} [get]
func (h *PartHandler) FindByID(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
