package vehiclehandler

import (
	"errors"
	"net/http"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/gin-gonic/gin"
)

// Create godoc
// @Summary     Cadastrar veiculo
// @Tags        vehicles
// @Accept      json
// @Produce     json
// @Param       body body commands.CreateVehicleRequest true "Dados do veiculo"
// @Success     201 {object} commands.VehicleResponse
// @Security    BearerAuth
// @Router      /vehicles [post]
func (h *VehicleHandler) Create(c *gin.Context) {
	var req commands.CreateVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.HandleBadRequest(c, err)
		return
	}

	vehicle, err := req.ToDomain()
	if err != nil {
		httputil.HandleBadRequest(c, err)
		return
	}
	if err := h.create.Execute(c.Request.Context(), vehicle); err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrNotFound):
			httputil.HandleError(c, err, msgCustomerNotFound)
		case errors.Is(err, domainerrors.ErrAlreadyExists):
			httputil.HandleError(c, err, msgAlreadyExists)
		default:
			httputil.HandleError(c, err)
		}
		return
	}

	c.JSON(http.StatusCreated, commands.ToVehicleResponse(vehicle))
}
