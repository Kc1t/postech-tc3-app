package parthandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Delete godoc
// @Summary     Deletar peca/insumo
// @Tags        parts
// @Produce     json
// @Param       id path string true "Part ID"
// @Security    BearerAuth
// @Router      /parts/{id} [delete]
func (h *PartHandler) Delete(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
