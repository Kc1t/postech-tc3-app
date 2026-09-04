package servicehandler

import (
	"net/http"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/gin-gonic/gin"
)

// Delete godoc
// @Summary     Deletar servico
// @Tags        services
// @Param       id path string true "Service ID"
// @Success     204
// @Security    BearerAuth
// @Router      /services/{id} [delete]
func (h *ServiceHandler) Delete(c *gin.Context) {
	if err := h.delete.Execute(c.Request.Context(), c.Param("id")); err != nil {
		httputil.HandleError(c, err, msgNotFound)
		return
	}
	c.Status(http.StatusNoContent)
}
