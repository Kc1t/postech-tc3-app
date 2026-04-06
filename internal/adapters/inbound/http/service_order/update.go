package serviceorderhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Update godoc
// @Summary     Atualizar ordem de servico
// @Tags        service-orders
// @Accept      json
// @Produce     json
// @Param       id path string true "ServiceOrder ID"
// @Security    BearerAuth
// @Router      /service-orders/{id} [put]
func (h *ServiceOrderHandler) Update(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
