package serviceorderhandler

import (
	"net/http"
	"strconv"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	"github.com/gin-gonic/gin"
)

// Reject godoc
// @Summary     Rejeitar ordem de servico (endpoint publico do cliente)
// @Tags        service-orders-customer
// @Accept      json
// @Param       code path int true "Codigo da OS"
// @Param       body body commands.ApproveRejectRequest true "CPF do cliente"
// @Success     204
// @Router      /service-orders/code/{code}/reject [post]
func (h *ServiceOrderHandler) Reject(c *gin.Context) {
	code, err := strconv.Atoi(c.Param("code"))
	if err != nil {
		httputil.HandleBadRequest(c, err)
		return
	}

	var req commands.ApproveRejectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.HandleBadRequest(c, err)
		return
	}

	if err := h.reject.Execute(c.Request.Context(), code, req.CustomerCPF); err != nil {
		httputil.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
