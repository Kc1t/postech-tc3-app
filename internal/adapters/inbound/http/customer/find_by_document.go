package customerhandler

import (
	"errors"
	"net/http"

	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/gin-gonic/gin"
)

// FindByDocument godoc
// @Summary     Buscar cliente por documento
// @Tags        customers
// @Produce     json
// @Param       document path string true "Customer Document"
// @Security    BearerAuth
// @Router      /customers/document/{document} [get]
func (h *CustomerHandler) FindByDocument(c *gin.Context) {
	customer, err := h.getByDocument.Execute(c.Request.Context(), c.Param("document"))
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
