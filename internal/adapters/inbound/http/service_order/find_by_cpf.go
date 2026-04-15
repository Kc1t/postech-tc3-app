package serviceorderhandler

import (
	"net/http"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	"github.com/gin-gonic/gin"
)

// FindByCPF godoc
// @Summary     Listar OS do cliente por CPF (endpoint publico do cliente)
// @Tags        service-orders-customer
// @Produce     json
// @Param       cpf query string true "CPF do cliente"
// @Success     200 {array} commands.ServiceOrderResponse
// @Router      /service-orders/customer [get]
func (h *ServiceOrderHandler) FindByCPF(c *gin.Context) {
	cpf := c.Query("cpf")
	if cpf == "" {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse{
			Status:  http.StatusText(http.StatusBadRequest),
			Message: "query parameter 'cpf' is required",
		})
		return
	}

	orders, err := h.listByCPF.Execute(c.Request.Context(), cpf)
	if err != nil {
		httputil.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, commands.ToServiceOrderListResponse(orders))
}
