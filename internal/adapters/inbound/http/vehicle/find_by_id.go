package vehiclehandler

import (
	"errors"
	"net/http"

	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
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
	vehicle, err := h.getByID.Execute(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, domainerrors.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "vehicle not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, commands.ToVehicleResponse(vehicle))
}
