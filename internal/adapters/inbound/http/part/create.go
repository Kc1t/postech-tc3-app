package parthandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create godoc
// @Summary     Cadastrar peca/insumo
// @Tags        parts
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Router      /parts [post]
func (h *PartHandler) Create(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
