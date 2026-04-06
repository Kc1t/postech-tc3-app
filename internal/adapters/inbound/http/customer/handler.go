package customerhandler

import (
	"github.com/fiap/postech-tc1/internal/ports"
	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	create        ports.CreateCustomerUseCase
	getByID       ports.GetCustomerUseCase
	getByDocument ports.GetCustomerByDocumentUseCase
	listAll       ports.ListCustomersUseCase
	update        ports.UpdateCustomerUseCase
	delete        ports.DeleteCustomerUseCase
}

func NewCustomerHandler(
	create ports.CreateCustomerUseCase,
	getByID ports.GetCustomerUseCase,
	getByDocument ports.GetCustomerByDocumentUseCase,
	listAll ports.ListCustomersUseCase,
	update ports.UpdateCustomerUseCase,
	delete ports.DeleteCustomerUseCase,
) *CustomerHandler {
	return &CustomerHandler{
		create:        create,
		getByID:       getByID,
		getByDocument: getByDocument,
		listAll:       listAll,
		update:        update,
		delete:        delete,
	}
}

func (h *CustomerHandler) SetupRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/customers")
	g.POST("", h.Create)
	g.GET("", h.FindAll)
	g.GET("/:id", h.FindByID)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}
