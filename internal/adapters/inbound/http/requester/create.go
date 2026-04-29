package requesterhandler

import (
	"net/http"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	"github.com/gin-gonic/gin"
)

// Create godoc
// @Summary     Criar cliente
// @Tags        requesters
// @Accept      json
// @Produce     json
// @Param       body body commands.CreateRequesterRequest true "Dados do cliente"
// @Success     201 {object} commands.RequesterResponse
// @Security    BearerAuth
// @Router      /requesters [post]
func (h *RequesterHandler) Create(c *gin.Context) {
	var req commands.CreateRequesterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.HandleBadRequest(c, err)
		return
	}

	requester, err := req.ToDomain()
	if err != nil {
		httputil.HandleBadRequest(c, err)
		return
	}
	if err := h.create.Execute(c.Request.Context(), requester); err != nil {
		httputil.HandleError(c, err, msgAlreadyExists)
		return
	}

	c.JSON(http.StatusCreated, commands.ToRequesterResponse(requester))
}
