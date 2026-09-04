package serviceorderhandler

import (
	"net/http"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	"github.com/gin-gonic/gin"
)

// FindByID godoc
// @Summary     Buscar ordem de servico por ID
// @Tags        service-orders
// @Produce     json
// @Param       id path string true "ServiceOrder ID"
// @Success     200 {object} commands.ServiceOrderResponse
// @Security    BearerAuth
// @Router      /service-orders/{id} [get]
func (h *ServiceOrderHandler) FindByID(c *gin.Context) {
	so, err := h.getByID.Execute(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.HandleError(c, err, msgNotFound)
		return
	}
	c.JSON(http.StatusOK, commands.ToServiceOrderResponse(so))
}
