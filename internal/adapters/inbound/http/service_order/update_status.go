package serviceorderhandler

import (
	"net/http"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/gin-gonic/gin"
)

// UpdateStatus godoc
// @Summary     Atualizar status da ordem de servico
// @Tags        service-orders
// @Accept      json
// @Produce     json
// @Param       id path string true "ServiceOrder ID"
// @Param       body body commands.UpdateStatusRequest true "Novo status (services e parts obrigatorios para awaiting_approval)"
// @Success     204
// @Security    BearerAuth
// @Router      /service-orders/{id}/status [put]
func (h *ServiceOrderHandler) UpdateStatus(c *gin.Context) {
	var req commands.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.HandleBadRequest(c, err)
		return
	}

	input := entities.StatusUpdate{
		ID:     c.Param("id"),
		Status: req.Status,
	}

	for _, s := range req.Services {
		input.ServiceCodes = append(input.ServiceCodes, s.Code)
	}
	for _, p := range req.Parts {
		input.Parts = append(input.Parts, entities.OrderPartItem{
			ManufacturerCode: p.ManufacturerCode,
			Quantity:         p.Quantity,
		})
	}

	if err := h.updateStatus.Execute(c.Request.Context(), input); err != nil {
		httputil.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
