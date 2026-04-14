package servicehandler

import (
	"net/http"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	"github.com/gin-gonic/gin"
)

// Create godoc
// @Summary     Criar servico
// @Tags        services
// @Accept      json
// @Produce     json
// @Param       body body commands.CreateServiceRequest true "Dados do servico"
// @Success     201 {object} commands.ServiceResponse
// @Security    BearerAuth
// @Router      /services [post]
func (h *ServiceHandler) Create(c *gin.Context) {
	var req commands.CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.HandleBadRequest(c, err)
		return
	}

	service := req.ToDomain()
	if err := h.create.Execute(c.Request.Context(), service); err != nil {
		httputil.HandleError(c, err, msgAlreadyExists)
		return
	}

	c.JSON(http.StatusCreated, commands.ToServiceResponse(service))
}
