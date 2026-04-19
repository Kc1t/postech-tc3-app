package serviceorderhandler

import (
	"github.com/fiap/postech-tc1/internal/ports"
	"github.com/gin-gonic/gin"
)

type ServiceOrderHandler struct {
	create       ports.CreateServiceOrderUseCase
	getByID      ports.GetServiceOrderUseCase
	listAll      ports.ListServiceOrdersUseCase
	updateStatus ports.UpdateServiceOrderStatusUseCase
	update       ports.UpdateServiceOrderUseCase
	delete       ports.DeleteServiceOrderUseCase
}

func NewServiceOrderHandler(
	create ports.CreateServiceOrderUseCase,
	getByID ports.GetServiceOrderUseCase,
	listAll ports.ListServiceOrdersUseCase,
	updateStatus ports.UpdateServiceOrderStatusUseCase,
	update ports.UpdateServiceOrderUseCase,
	delete ports.DeleteServiceOrderUseCase,
) *ServiceOrderHandler {
	return &ServiceOrderHandler{
		create:       create,
		getByID:      getByID,
		listAll:      listAll,
		updateStatus: updateStatus,
		update:       update,
		delete:       delete,
	}
}

func (h *ServiceOrderHandler) SetupRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/service-orders")
	g.POST("", h.Create)
	g.GET("", h.FindAll)
	g.GET("/:id", h.FindByID)
	g.PUT("/:id/status", h.UpdateStatus)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}
