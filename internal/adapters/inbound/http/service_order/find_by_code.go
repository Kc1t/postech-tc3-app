package serviceorderhandler

import (
	"fmt"
	"net/http"
	"strconv"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	"github.com/gin-gonic/gin"
)

// FindByCode godoc
// @Summary     Consultar OS pelo codigo (endpoint publico do cliente)
// @Tags        service-orders-customer
// @Produce     json
// @Param       code path int true "Codigo da OS"
// @Param       cpf query string true "CPF do cliente"
// @Success     200 {object} commands.ServiceOrderResponse
// @Router      /service-orders/code/{code} [get]
func (h *ServiceOrderHandler) FindByCode(c *gin.Context) {
	code, err := strconv.Atoi(c.Param("code"))
	if err != nil {
		httputil.HandleBadRequest(c, err)
		return
	}

	cpf := c.Query("cpf")
	if cpf == "" {
		httputil.HandleBadRequest(c, fmt.Errorf("query parameter 'cpf' is required"))
		return
	}

	so, err := h.getByCode.Execute(c.Request.Context(), code, cpf)
	if err != nil {
		httputil.HandleError(c, err, msgNotFound)
		return
	}

	c.JSON(http.StatusOK, commands.ToServiceOrderResponse(so))
}
