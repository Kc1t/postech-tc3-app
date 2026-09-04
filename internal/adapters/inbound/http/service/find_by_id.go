package servicehandler

import (
	"net/http"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
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
		httputil.HandleError(c, err, msgNotFound)
		return
	}
	c.JSON(http.StatusOK, commands.ToServiceResponse(service))
}
