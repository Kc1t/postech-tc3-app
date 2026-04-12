package servicehandler

import (
	"errors"
	"net/http"

	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/gin-gonic/gin"
)

// FindByID godoc
// @Summary     Buscar servico por ID
// @Tags        services
// @Produce     json
// @Param       id path string true "Service ID"
// @Success     200 {object} commands.ServiceResponse
// @Security    BearerAuth
// @Router      /services/{id} [get]
func (h *ServiceHandler) FindByID(c *gin.Context) {
	service, err := h.getByID.Execute(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, domainerrors.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, commands.ToServiceResponse(service))
}
