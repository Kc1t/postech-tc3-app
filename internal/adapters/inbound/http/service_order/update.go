package serviceorderhandler

import (
	"net/http"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	"github.com/gin-gonic/gin"
)

// Update godoc
// @Summary     Atualizar ordem de servico
// @Tags        service-orders
// @Accept      json
// @Produce     json
// @Param       id path string true "ServiceOrder ID"
// @Param       body body commands.UpdateServiceOrderRequest true "Dados para atualizar"
// @Success     200 {object} commands.ServiceOrderResponse
// @Security    BearerAuth
// @Router      /service-orders/{id} [put]
func (h *ServiceOrderHandler) Update(c *gin.Context) {
	var req commands.UpdateServiceOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.HandleBadRequest(c, err)
		return
	}

	so, err := h.getByID.Execute(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.HandleError(c, err, msgNotFound)
		return
	}

	if req.Notes != nil {
		so.SetNotes(*req.Notes)
	}

	if err := h.update.Execute(c.Request.Context(), so); err != nil {
		httputil.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, commands.ToServiceOrderResponse(so))
}
