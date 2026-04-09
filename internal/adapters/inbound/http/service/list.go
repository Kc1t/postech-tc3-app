package servicehandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// FindAll godoc
// @Summary     Listar servicos
// @Tags        services
// @Produce     json
// @Security    BearerAuth
// @Router      /services [get]
func (h *ServiceHandler) FindAll(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
