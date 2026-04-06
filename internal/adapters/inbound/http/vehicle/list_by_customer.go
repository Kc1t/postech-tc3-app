package vehiclehandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListByCustomer godoc
// @Summary     Listar veiculos por cliente
// @Tags        vehicles
// @Produce     json
// @Param       customer_id path string true "Customer ID"
// @Security    BearerAuth
// @Router      /customers/{customer_id}/vehicles [get]
func (h *VehicleHandler) ListByCustomer(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
