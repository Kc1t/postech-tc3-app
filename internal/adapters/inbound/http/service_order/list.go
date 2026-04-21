package serviceorderhandler

import (
	"net/http"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	"github.com/gin-gonic/gin"
)

// FindAll godoc
// @Summary     Listar ordens de servico
// @Tags        service-orders
// @Produce     json
// @Success     200 {array} commands.ServiceOrderResponse
// @Security    BearerAuth
// @Router      /service-orders [get]
func (h *ServiceOrderHandler) FindAll(c *gin.Context) {
	orders, err := h.listAll.Execute(c.Request.Context())
	if err != nil {
		httputil.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, commands.ToServiceOrderListResponse(orders))
}
