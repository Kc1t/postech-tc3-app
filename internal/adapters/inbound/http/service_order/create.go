package serviceorderhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create godoc
// @Summary     Criar ordem de servico
// @Tags        service-orders
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Router      /service-orders [post]
func (h *ServiceOrderHandler) Create(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
