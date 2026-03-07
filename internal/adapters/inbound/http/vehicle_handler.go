package handler

import (
	"net/http"

	"github.com/fiap/postech-tc1/internal/ports"
	"github.com/gin-gonic/gin"
)

type VehicleHandler struct {
	useCase ports.VehicleUseCase
}

func NewVehicleHandler(uc ports.VehicleUseCase) *VehicleHandler {
	return &VehicleHandler{useCase: uc}
}

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

// FindAll godoc
// @Summary     Listar veiculos
// @Tags        vehicles
// @Produce     json
// @Security    BearerAuth
// @Router      /vehicles [get]
func (h *VehicleHandler) FindAll(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

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
