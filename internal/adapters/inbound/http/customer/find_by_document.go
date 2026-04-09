package customerhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// FindByDocument godoc
// @Summary     Buscar cliente por documento
// @Tags        customers
// @Produce     json
// @Param       document path string true "Customer Document"
// @Security    BearerAuth
// @Router      /customers/document/{document} [get]
func (h *CustomerHandler) FindByDocument(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
