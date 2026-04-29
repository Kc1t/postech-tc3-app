package requesterhandler

import (
	"net/http"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/gin-gonic/gin"
)

// Delete godoc
// @Summary     Deletar cliente
// @Tags        requesters
// @Param       id path string true "Requester ID"
// @Success     204
// @Security    BearerAuth
// @Router      /requesters/{id} [delete]
func (h *RequesterHandler) Delete(c *gin.Context) {
	if err := h.delete.Execute(c.Request.Context(), c.Param("id")); err != nil {
		httputil.HandleError(c, err, msgNotFound)
		return
	}
	c.Status(http.StatusNoContent)
}
