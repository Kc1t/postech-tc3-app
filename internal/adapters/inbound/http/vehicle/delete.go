package vehiclehandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Delete godoc
// @Summary     Deletar veiculo
// @Tags        vehicles
// @Produce     json
// @Param       id path string true "Vehicle ID"
// @Security    BearerAuth
// @Router      /vehicles/{id} [delete]
func (h *VehicleHandler) Delete(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
