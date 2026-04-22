package customerhandler

import (
	"github.com/fiap/postech-tc1/internal/ports"
	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	create           ports.CreateCustomerUseCase
	getByID          ports.GetCustomerUseCase
	getByDocument    ports.GetCustomerByDocumentUseCase
	listAll          ports.ListCustomersUseCase
	update           ports.UpdateCustomerUseCase
	delete           ports.DeleteCustomerUseCase
	listVehicles     ports.ListVehiclesByCustomerUseCase
}

func NewCustomerHandler(
	create ports.CreateCustomerUseCase,
	getByID ports.GetCustomerUseCase,
	getByDocument ports.GetCustomerByDocumentUseCase,
	listAll ports.ListCustomersUseCase,
	update ports.UpdateCustomerUseCase,
	delete ports.DeleteCustomerUseCase,
	listVehicles ports.ListVehiclesByCustomerUseCase,
) *CustomerHandler {
	return &CustomerHandler{
		create:        create,
		getByID:       getByID,
		getByDocument: getByDocument,
		listAll:       listAll,
		update:        update,
		delete:        delete,
		listVehicles:  listVehicles,
	}
}

func (h *CustomerHandler) SetupRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/customers")
	g.POST("", h.Create)
	g.GET("", h.FindAll)
	g.GET("/document/:document", h.FindByDocument)
	g.GET("/:id", h.FindByID)
	g.GET("/:id/vehicles", h.ListVehicles)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}
