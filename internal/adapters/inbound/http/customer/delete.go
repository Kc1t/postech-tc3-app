package customerhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Delete godoc
// @Summary     Deletar cliente
// @Tags        customers
// @Produce     json
// @Param       id path string true "Customer ID"
// @Security    BearerAuth
// @Router      /customers/{id} [delete]
func (h *CustomerHandler) Delete(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
