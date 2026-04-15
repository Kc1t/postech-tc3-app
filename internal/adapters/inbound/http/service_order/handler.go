package serviceorderhandler

import (
	"github.com/fiap/postech-tc1/internal/ports"
	"github.com/gin-gonic/gin"
)

type ServiceOrderHandler struct {
	create       ports.CreateServiceOrderUseCase
	getByID      ports.GetServiceOrderUseCase
	getByCode    ports.GetServiceOrderByCodeUseCase
	listAll      ports.ListServiceOrdersUseCase
	listByCPF    ports.ListServiceOrdersByCPFUseCase
	updateStatus ports.UpdateServiceOrderStatusUseCase
	update       ports.UpdateServiceOrderUseCase
	delete       ports.DeleteServiceOrderUseCase
	approve      ports.ApproveServiceOrderUseCase
	reject       ports.RejectServiceOrderUseCase
}

func NewServiceOrderHandler(
	create ports.CreateServiceOrderUseCase,
	getByID ports.GetServiceOrderUseCase,
	getByCode ports.GetServiceOrderByCodeUseCase,
	listAll ports.ListServiceOrdersUseCase,
	listByCPF ports.ListServiceOrdersByCPFUseCase,
	updateStatus ports.UpdateServiceOrderStatusUseCase,
	update ports.UpdateServiceOrderUseCase,
	delete ports.DeleteServiceOrderUseCase,
	approve ports.ApproveServiceOrderUseCase,
	reject ports.RejectServiceOrderUseCase,
) *ServiceOrderHandler {
	return &ServiceOrderHandler{
		create:       create,
		getByID:      getByID,
		getByCode:    getByCode,
		listAll:      listAll,
		listByCPF:    listByCPF,
		updateStatus: updateStatus,
		update:       update,
		delete:       delete,
		approve:      approve,
		reject:       reject,
	}
}

// SetupRoutes registra as rotas administrativas (protegidas por JWT).
func (h *ServiceOrderHandler) SetupRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/service-orders")
	g.POST("", h.Create)
	g.GET("", h.FindAll)
	g.GET("/:id", h.FindByID)
	g.PUT("/:id/status", h.UpdateStatus)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}

// SetupPublicRoutes registra as rotas publicas do cliente (sem JWT).
func (h *ServiceOrderHandler) SetupPublicRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/service-orders")
	g.GET("/code/:code", h.FindByCode)
	g.GET("/customer", h.FindByCPF)
	g.POST("/code/:code/approve", h.Approve)
	g.POST("/code/:code/reject", h.Reject)
}
