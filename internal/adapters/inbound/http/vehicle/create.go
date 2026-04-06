package vehiclehandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create godoc
// @Summary     Cadastrar veiculo
// @Tags        vehicles
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Router      /vehicles [post]
func (h *VehicleHandler) Create(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
