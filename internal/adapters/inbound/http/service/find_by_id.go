package servicehandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// FindByID godoc
// @Summary     Buscar servico por ID
// @Tags        services
// @Produce     json
// @Param       id path string true "Service ID"
// @Security    BearerAuth
// @Router      /services/{id} [get]
func (h *ServiceHandler) FindByID(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
