package servicehandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Update godoc
// @Summary     Atualizar servico
// @Tags        services
// @Accept      json
// @Produce     json
// @Param       id path string true "Service ID"
// @Security    BearerAuth
// @Router      /services/{id} [put]
func (h *ServiceHandler) Update(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
