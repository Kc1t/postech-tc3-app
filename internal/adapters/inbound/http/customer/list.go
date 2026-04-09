package customerhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// FindAll godoc
// @Summary     Listar clientes
// @Tags        customers
// @Produce     json
// @Security    BearerAuth
// @Router      /customers [get]
func (h *CustomerHandler) FindAll(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
