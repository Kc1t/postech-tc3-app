package serviceorderhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// FindByID godoc
// @Summary     Buscar ordem de servico por ID
// @Tags        service-orders
// @Produce     json
// @Param       id path string true "ServiceOrder ID"
// @Security    BearerAuth
// @Router      /service-orders/{id} [get]
func (h *ServiceOrderHandler) FindByID(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
