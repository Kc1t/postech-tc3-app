package parthandler

import (
	"net/http"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/gin-gonic/gin"
)

// Delete godoc
// @Summary     Deletar peca/insumo
// @Tags        parts
// @Param       id path string true "Part ID"
// @Success     204
// @Security    BearerAuth
// @Router      /parts/{id} [delete]
func (h *PartHandler) Delete(c *gin.Context) {
	if err := h.delete.Execute(c.Request.Context(), c.Param("id")); err != nil {
		httputil.HandleError(c, err, msgNotFound)
		return
	}
	c.Status(http.StatusNoContent)
}
