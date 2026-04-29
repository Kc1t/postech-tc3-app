package serviceorderhandler

import (
	"net/http"
	"strconv"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/gin-gonic/gin"
)

// UpdateStatusByCode godoc
// @Summary     Alterar status da OS pelo codigo (endpoint publico do cliente)
// @Tags        service-orders-requester
// @Accept      json
// @Produce     json
// @Param       code path int true "Codigo da OS"
// @Param       body body commands.UpdateStatusByCodeRequest true "Novo status e documento do cliente"
// @Success     204
// @Router      /service-orders/code/{code}/status [put]
func (h *ServiceOrderHandler) UpdateStatusByCode(c *gin.Context) {
	code, err := strconv.Atoi(c.Param("code"))
	if err != nil {
		httputil.HandleBadRequest(c, err)
		return
	}

	var req commands.UpdateStatusByCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.HandleBadRequest(c, err)
		return
	}

	doc, err := entities.NewDocument(req.RequesterDocument)
	if err != nil {
		httputil.HandleBadRequest(c, err)
		return
	}

	if err := h.updateStatusByCode.Execute(c.Request.Context(), code, doc.Value(), req.Status); err != nil {
		httputil.HandleError(c, err, msgNotFound)
		return
	}

	c.Status(http.StatusNoContent)
}
