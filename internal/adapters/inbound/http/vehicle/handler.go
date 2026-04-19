package vehiclehandler

import (
	"github.com/fiap/postech-tc1/internal/ports"
	"github.com/gin-gonic/gin"
)

type VehicleHandler struct {
	create  ports.CreateVehicleUseCase
	getByID ports.GetVehicleUseCase
	listAll ports.ListVehiclesUseCase
	update  ports.UpdateVehicleUseCase
	delete  ports.DeleteVehicleUseCase
}

func NewVehicleHandler(
	create ports.CreateVehicleUseCase,
	getByID ports.GetVehicleUseCase,
	listAll ports.ListVehiclesUseCase,
	update ports.UpdateVehicleUseCase,
	delete ports.DeleteVehicleUseCase,
) *VehicleHandler {
	return &VehicleHandler{
		create:  create,
		getByID: getByID,
		listAll: listAll,
		update:  update,
		delete:  delete,
	}
}

func (h *VehicleHandler) SetupRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/vehicles")
	g.POST("", h.Create)
	g.GET("", h.FindAll)
	g.GET("/:id", h.FindByID)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}
