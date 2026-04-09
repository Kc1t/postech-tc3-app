package vehiclehandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// FindByID godoc
// @Summary     Buscar veiculo por ID
// @Tags        vehicles
// @Produce     json
// @Param       id path string true "Vehicle ID"
// @Security    BearerAuth
// @Router      /vehicles/{id} [get]
func (h *VehicleHandler) FindByID(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
