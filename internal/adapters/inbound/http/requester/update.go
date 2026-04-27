package requesterhandler

import (
	"net/http"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	"github.com/gin-gonic/gin"
)

// Update godoc
// @Summary     Atualizar cliente
// @Tags        requesters
// @Accept      json
// @Produce     json
// @Param       id path string true "Requester ID"
// @Param       body body commands.UpdateRequesterRequest true "Dados para atualizar"
// @Success     200 {object} commands.RequesterResponse
// @Security    BearerAuth
// @Router      /requesters/{id} [put]
func (h *RequesterHandler) Update(c *gin.Context) {
	var req commands.UpdateRequesterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.HandleBadRequest(c, err)
		return
	}

	requester, err := h.getByID.Execute(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.HandleError(c, err, msgNotFound)
		return
	}

	if req.Name != "" {
		requester.SetName(req.Name)
	}
	if req.Email != "" {
		requester.SetEmail(req.Email)
	}
	if req.Phone != "" {
		requester.SetPhone(req.Phone)
	}

	if err := h.update.Execute(c.Request.Context(), requester); err != nil {
		httputil.HandleError(c, err)
		return
	}


	c.JSON(http.StatusOK, commands.ToRequesterResponse(requester))
}
