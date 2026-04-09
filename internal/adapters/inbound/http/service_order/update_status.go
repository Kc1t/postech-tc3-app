package serviceorderhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// UpdateStatus godoc
// @Summary     Atualizar status da ordem de servico
// @Tags        service-orders
// @Accept      json
// @Produce     json
// @Param       id path string true "ServiceOrder ID"
// @Security    BearerAuth
// @Router      /service-orders/{id}/status [put]
func (h *ServiceOrderHandler) UpdateStatus(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
