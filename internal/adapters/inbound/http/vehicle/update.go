package vehiclehandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Update godoc
// @Summary     Atualizar veiculo
// @Tags        vehicles
// @Accept      json
// @Produce     json
// @Param       id path string true "Vehicle ID"
// @Security    BearerAuth
// @Router      /vehicles/{id} [put]
func (h *VehicleHandler) Update(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
