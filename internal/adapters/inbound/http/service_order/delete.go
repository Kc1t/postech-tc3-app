package serviceorderhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Delete godoc
// @Summary     Deletar ordem de servico
// @Tags        service-orders
// @Produce     json
// @Param       id path string true "ServiceOrder ID"
// @Security    BearerAuth
// @Router      /service-orders/{id} [delete]
func (h *ServiceOrderHandler) Delete(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
