package customerhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create godoc
// @Summary     Criar cliente
// @Tags        customers
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Router      /customers [post]
func (h *CustomerHandler) Create(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
