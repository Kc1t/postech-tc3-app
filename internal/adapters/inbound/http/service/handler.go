package servicehandler

import (
	"github.com/fiap/postech-tc1/internal/ports"
	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	create  ports.CreateServiceUseCase
	getByID ports.GetServiceUseCase
	listAll ports.ListServicesUseCase
	update  ports.UpdateServiceUseCase
	delete  ports.DeleteServiceUseCase
}

func NewServiceHandler(
	create ports.CreateServiceUseCase,
	getByID ports.GetServiceUseCase,
	listAll ports.ListServicesUseCase,
	update ports.UpdateServiceUseCase,
	delete ports.DeleteServiceUseCase,
) *ServiceHandler {
	return &ServiceHandler{
		create:  create,
		getByID: getByID,
		listAll: listAll,
		update:  update,
		delete:  delete,
	}
}

func (h *ServiceHandler) SetupRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/services")
	g.POST("", h.Create)
	g.GET("", h.FindAll)
	g.GET("/:id", h.FindByID)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}
