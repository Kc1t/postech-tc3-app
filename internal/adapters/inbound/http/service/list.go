package servicehandler

import (
	"net/http"

	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	"github.com/gin-gonic/gin"
)

// FindAll godoc
// @Summary     Listar servicos
// @Tags        services
// @Produce     json
// @Success     200 {array} commands.ServiceResponse
// @Security    BearerAuth
// @Router      /services [get]
func (h *ServiceHandler) FindAll(c *gin.Context) {
	services, err := h.listAll.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, commands.ToServiceListResponse(services))
}
