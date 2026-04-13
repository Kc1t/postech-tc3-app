package vehiclehandler

import (
	"net/http"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	"github.com/gin-gonic/gin"
)

// Update godoc
// @Summary     Atualizar veiculo
// @Tags        vehicles
// @Accept      json
// @Produce     json
// @Param       id path string true "Vehicle ID"
// @Param       body body commands.UpdateVehicleRequest true "Dados para atualizar"
// @Success     200 {object} commands.VehicleResponse
// @Security    BearerAuth
// @Router      /vehicles/{id} [put]
func (h *VehicleHandler) Update(c *gin.Context) {
	var req commands.UpdateVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.HandleBadRequest(c, err)
		return
	}

	vehicle, err := h.getByID.Execute(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.HandleError(c, err, msgNotFound)
		return
	}

	if req.Plate != "" {
		vehicle.SetPlate(req.Plate)
	}
	if req.Brand != "" {
		vehicle.SetBrand(req.Brand)
	}
	if req.Model != "" {
		vehicle.SetModel(req.Model)
	}
	if req.Year != 0 {
		vehicle.SetYear(req.Year)
	}

	if err := h.update.Execute(c.Request.Context(), vehicle); err != nil {
		httputil.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, commands.ToVehicleResponse(vehicle))
}
