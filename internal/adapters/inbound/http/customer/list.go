package customerhandler

import (
	"net/http"

	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	"github.com/gin-gonic/gin"
)

// FindAll godoc
// @Summary     Listar clientes
// @Tags        customers
// @Produce     json
// @Success     200 {array} commands.CustomerResponse
// @Security    BearerAuth
// @Router      /customers [get]
func (h *CustomerHandler) FindAll(c *gin.Context) {
	customers, err := h.listAll.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, commands.ToCustomerListResponse(customers))
}
