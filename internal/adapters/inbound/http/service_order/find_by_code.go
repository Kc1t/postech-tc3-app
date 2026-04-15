package serviceorderhandler

import (
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
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse{
			Status:  http.StatusText(http.StatusBadRequest),
			Message: "query parameter 'cpf' is required",
		})
		return
	}

	so, err := h.getByCode.Execute(c.Request.Context(), code, cpf)
	if err != nil {
		httputil.HandleError(c, err, msgNotFound)
		return
	}

	c.JSON(http.StatusOK, commands.ToServiceOrderResponse(so))
}
