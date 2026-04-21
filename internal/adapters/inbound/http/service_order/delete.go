package serviceorderhandler

import (
	"net/http"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/gin-gonic/gin"
)

// Delete godoc
// @Summary     Deletar ordem de servico
// @Tags        service-orders
// @Produce     json
// @Param       id path string true "ServiceOrder ID"
// @Success     204
// @Security    BearerAuth
// @Router      /service-orders/{id} [delete]
func (h *ServiceOrderHandler) Delete(c *gin.Context) {
	if err := h.delete.Execute(c.Request.Context(), c.Param("id")); err != nil {
		httputil.HandleError(c, err, msgNotFound)
		return
	}
	c.Status(http.StatusNoContent)
}
