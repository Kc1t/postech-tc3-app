package serviceorderhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// FindAll godoc
// @Summary     Listar ordens de servico
// @Tags        service-orders
// @Produce     json
// @Security    BearerAuth
// @Router      /service-orders [get]
func (h *ServiceOrderHandler) FindAll(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
