package servicehandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create godoc
// @Summary     Criar servico
// @Tags        services
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Router      /services [post]
func (h *ServiceHandler) Create(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
