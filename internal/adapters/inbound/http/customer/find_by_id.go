package customerhandler

import (
	"net/http"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	"github.com/gin-gonic/gin"
)

// FindByID godoc
// @Summary     Buscar cliente por ID
// @Tags        customers
// @Produce     json
// @Param       id path string true "Customer ID"
// @Success     200 {object} commands.CustomerResponse
// @Security    BearerAuth
// @Router      /customers/{id} [get]
func (h *CustomerHandler) FindByID(c *gin.Context) {
	customer, err := h.getByID.Execute(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.HandleError(c, err, msgNotFound)
		return
	}
	c.JSON(http.StatusOK, commands.ToCustomerResponse(customer))
}
