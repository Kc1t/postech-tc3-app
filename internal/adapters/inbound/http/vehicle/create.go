package vehiclehandler

import (
	"errors"
	"net/http"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vehicle := req.ToDomain()
	if err := h.create.Execute(c.Request.Context(), req.CustomerDocument, vehicle); err != nil {
		if errors.Is(err, domainerrors.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, commands.ToVehicleResponse(vehicle))
}
