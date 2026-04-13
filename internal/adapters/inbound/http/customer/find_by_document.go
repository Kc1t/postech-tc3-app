package customerhandler

import (
	"net/http"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
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
		httputil.HandleError(c, err, msgNotFound)
		return
	}
	c.JSON(http.StatusOK, commands.ToCustomerResponse(customer))
}
