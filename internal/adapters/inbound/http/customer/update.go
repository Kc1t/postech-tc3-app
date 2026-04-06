package customerhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Update godoc
// @Summary     Atualizar cliente
// @Tags        customers
// @Accept      json
// @Produce     json
// @Param       id path string true "Customer ID"
// @Security    BearerAuth
// @Router      /customers/{id} [put]
func (h *CustomerHandler) Update(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
