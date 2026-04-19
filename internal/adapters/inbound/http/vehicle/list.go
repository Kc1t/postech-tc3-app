package vehiclehandler

import (
	"net/http"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	"github.com/gin-gonic/gin"
)

// FindAll godoc
// @Summary     Listar veiculos
// @Tags        vehicles
// @Produce     json
// @Success     200 {array} commands.VehicleResponse
// @Security    BearerAuth
// @Router      /vehicles [get]
func (h *VehicleHandler) FindAll(c *gin.Context) {
	vehicles, err := h.listAll.Execute(c.Request.Context())
	if err != nil {
		httputil.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, commands.ToVehicleListResponse(vehicles))
}
