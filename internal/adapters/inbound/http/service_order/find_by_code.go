package serviceorderhandler

import (
	"fmt"
	"net/http"
	"strconv"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/gin-gonic/gin"
)

// FindByCode godoc
// @Summary     Consultar OS pelo codigo (endpoint publico do cliente)
// @Tags        service-orders-requester
// @Produce     json
// @Param       code path int true "Codigo da OS"
// @Param       document query string true "CPF ou CNPJ do cliente"
// @Success     200 {object} commands.ServiceOrderResponse
// @Router      /service-orders/code/{code} [get]
func (h *ServiceOrderHandler) FindByCode(c *gin.Context) {
	code, err := strconv.Atoi(c.Param("code"))
	if err != nil {
		httputil.HandleBadRequest(c, err)
		return
	}

	raw := c.Query("document")
	if raw == "" {
		httputil.HandleBadRequest(c, fmt.Errorf("query parameter 'document' is required"))
		return
	}

	doc, err := entities.NewDocument(raw)
	if err != nil {
		httputil.HandleBadRequest(c, err)
		return
	}

	so, err := h.getByCode.Execute(c.Request.Context(), code, doc.Value())
	if err != nil {
		httputil.HandleError(c, err, msgNotFound)
		return
	}

	c.JSON(http.StatusOK, commands.ToServiceOrderResponse(so))
}
