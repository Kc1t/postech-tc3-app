package serviceorderhandler

import (
	"github.com/fiap/postech-tc1/internal/ports"
	"github.com/gin-gonic/gin"
)

type ServiceOrderHandler struct {
	create               ports.CreateServiceOrderUseCase
	getByID              ports.GetServiceOrderUseCase
	getByCode            ports.GetServiceOrderByCodeUseCase
	listAll              ports.ListServiceOrdersUseCase
	listByDocument       ports.ListServiceOrdersByDocumentUseCase
	updateStatus         ports.UpdateServiceOrderStatusUseCase
	update               ports.UpdateServiceOrderUseCase
	delete               ports.DeleteServiceOrderUseCase
	updateStatusByCode   ports.UpdateServiceOrderStatusByCodeUseCase
	averageExecutionTime ports.GetAverageExecutionTimeUseCase
}

func NewServiceOrderHandler(
	create ports.CreateServiceOrderUseCase,
	getByID ports.GetServiceOrderUseCase,
	getByCode ports.GetServiceOrderByCodeUseCase,
	listAll ports.ListServiceOrdersUseCase,
	listByDocument ports.ListServiceOrdersByDocumentUseCase,
	updateStatus ports.UpdateServiceOrderStatusUseCase,
	update ports.UpdateServiceOrderUseCase,
	delete ports.DeleteServiceOrderUseCase,
	updateStatusByCode ports.UpdateServiceOrderStatusByCodeUseCase,
	averageExecutionTime ports.GetAverageExecutionTimeUseCase,
) *ServiceOrderHandler {
	return &ServiceOrderHandler{
		create:               create,
		getByID:              getByID,
		getByCode:            getByCode,
		listAll:              listAll,
		listByDocument:       listByDocument,
		updateStatus:         updateStatus,
		update:               update,
		delete:               delete,
		updateStatusByCode:   updateStatusByCode,
		averageExecutionTime: averageExecutionTime,
	}
}

// SetupRoutes registra as rotas administrativas (protegidas por JWT).
func (h *ServiceOrderHandler) SetupRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/service-orders")
	g.POST("", h.Create)
	g.GET("", h.FindAll)
	g.GET("/metrics/execution-time", h.AverageExecutionTime)
	g.GET("/:id", h.FindByID)
	g.PUT("/:id/status", h.UpdateStatus)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}

// SetupPublicRoutes registra as rotas publicas do cliente (sem JWT).
func (h *ServiceOrderHandler) SetupPublicRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/service-orders")
	g.GET("/code/:code", h.FindByCode)
	g.GET("/requester", h.FindByDocument)
	g.PUT("/code/:code/status", h.UpdateStatusByCode)
}
