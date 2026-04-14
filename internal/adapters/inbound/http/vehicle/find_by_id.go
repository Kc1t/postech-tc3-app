package vehiclehandler

import (
	"net/http"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
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
		httputil.HandleError(c, err, msgNotFound)
		return
	}
	c.JSON(http.StatusOK, commands.ToVehicleResponse(vehicle))
}
