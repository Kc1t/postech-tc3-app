package servicehandler

import (
	"net/http"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service := req.ToDomain()
	if err := h.create.Execute(c.Request.Context(), service); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, commands.ToServiceResponse(service))
}
