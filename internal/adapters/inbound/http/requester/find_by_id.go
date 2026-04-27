package requesterhandler

import (
	"net/http"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	"github.com/gin-gonic/gin"
)

// FindByID godoc
// @Summary     Buscar cliente por ID
// @Tags        requesters
// @Produce     json
// @Param       id path string true "Requester ID"
// @Success     200 {object} commands.RequesterResponse
// @Security    BearerAuth
// @Router      /requesters/{id} [get]
func (h *RequesterHandler) FindByID(c *gin.Context) {
	requester, err := h.getByID.Execute(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.HandleError(c, err, msgNotFound)
		return
	}
	c.JSON(http.StatusOK, commands.ToRequesterResponse(requester))
}
