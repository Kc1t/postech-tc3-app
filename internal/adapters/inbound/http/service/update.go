package servicehandler

import (
	"net/http"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	"github.com/gin-gonic/gin"
)

// Update godoc
// @Summary     Atualizar servico
// @Tags        services
// @Accept      json
// @Produce     json
// @Param       id path string true "Service ID"
// @Param       body body commands.UpdateServiceRequest true "Dados para atualizar"
// @Success     200 {object} commands.ServiceResponse
// @Security    BearerAuth
// @Router      /services/{id} [put]
func (h *ServiceHandler) Update(c *gin.Context) {
	var req commands.UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.HandleBadRequest(c, err)
		return
	}

	service, err := h.getByID.Execute(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.HandleError(c, err, msgNotFound)
		return
	}

	if req.Code > 0 {
		service.SetCode(req.Code)
	}
	if req.Name != "" {
		service.SetName(req.Name)
	}
	if req.Description != "" {
		service.SetDescription(req.Description)
	}
	if req.Price > 0 {
		service.SetPrice(req.Price)
	}
	if req.DurationMin > 0 {
		service.SetDurationMin(req.DurationMin)
	}

	if err := h.update.Execute(c.Request.Context(), service); err != nil {
		httputil.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, commands.ToServiceResponse(service))
}
