package vehiclehandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// FindAll godoc
// @Summary     Listar veiculos
// @Tags        vehicles
// @Produce     json
// @Security    BearerAuth
// @Router      /vehicles [get]
func (h *VehicleHandler) FindAll(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
