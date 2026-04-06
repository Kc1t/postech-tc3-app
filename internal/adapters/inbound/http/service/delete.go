package servicehandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Delete godoc
// @Summary     Deletar servico
// @Tags        services
// @Produce     json
// @Param       id path string true "Service ID"
// @Security    BearerAuth
// @Router      /services/{id} [delete]
func (h *ServiceHandler) Delete(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
