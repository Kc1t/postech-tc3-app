package requesterhandler

import (
	"net/http"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	"github.com/gin-gonic/gin"
)

// ListVehicles godoc
// @Summary     Listar veiculos de um cliente
// @Description Retorna todos os veiculos cadastrados para o cliente identificado pelo ID.
// @Tags        requesters
// @Produce     json
// @Param       id path string true "ID do cliente (UUID)"
// @Success     200 {array} commands.VehicleResponse
// @Failure     500 {object} httputil.ErrorResponse
// @Security    BearerAuth
// @Router      /requesters/{id}/vehicles [get]
func (h *RequesterHandler) ListVehicles(c *gin.Context) {
	vehicles, err := h.listVehicles.Execute(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, commands.ToVehicleListResponse(vehicles))
}
