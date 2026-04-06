package customerhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// FindByID godoc
// @Summary     Buscar cliente por ID
// @Tags        customers
// @Produce     json
// @Param       id path string true "Customer ID"
// @Security    BearerAuth
// @Router      /customers/{id} [get]
func (h *CustomerHandler) FindByID(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
