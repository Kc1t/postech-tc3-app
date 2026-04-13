package customerhandler

import (
	"errors"
	"net/http"

	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
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
		if errors.Is(err, domainerrors.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, commands.ToCustomerResponse(customer))
}
