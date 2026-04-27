package serviceorderhandler

import (
	"net/http"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/gin-gonic/gin"
)

// Create godoc
// @Summary     Criar ordem de servico
// @Tags        service-orders
// @Accept      json
// @Produce     json
// @Param       body body commands.CreateServiceOrderRequest true "Dados da OS"
// @Success     201 {object} commands.ServiceOrderResponse
// @Security    BearerAuth
// @Router      /service-orders [post]
func (h *ServiceOrderHandler) Create(c *gin.Context) {
	var req commands.CreateServiceOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.HandleBadRequest(c, err)
		return
	}

	doc, err := entities.NewDocument(req.RequesterDocument)
	if err != nil {
		httputil.HandleBadRequest(c, err)
		return
	}

	input := entities.ServiceOrderInput{
		RequesterDocument: doc.Value(),
		VehiclePlate:     req.VehiclePlate,
		Notes:            req.Notes,
	}

	so, err := h.create.Execute(c.Request.Context(), input)
	if err != nil {
		httputil.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, commands.ToServiceOrderResponse(so))
}
