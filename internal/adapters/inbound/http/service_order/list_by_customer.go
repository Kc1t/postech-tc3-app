package serviceorderhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListByCustomer godoc
// @Summary     Listar ordens de servico por cliente
// @Tags        service-orders
// @Produce     json
// @Param       id path string true "Customer ID"
// @Security    BearerAuth
// @Router      /customers/{id}/service-orders [get]
func (h *ServiceOrderHandler) ListByCustomer(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
