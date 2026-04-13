package parthandler

import (
	"net/http"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	"github.com/gin-gonic/gin"
)

// Update godoc
// @Summary     Atualizar peca/insumo
// @Tags        parts
// @Accept      json
// @Produce     json
// @Param       id path string true "Part ID"
// @Param       body body commands.UpdatePartRequest true "Dados para atualizar"
// @Success     200 {object} commands.PartResponse
// @Security    BearerAuth
// @Router      /parts/{id} [put]
func (h *PartHandler) Update(c *gin.Context) {
	var req commands.UpdatePartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.HandleBadRequest(c, err)
		return
	}

	part, err := h.getByID.Execute(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.HandleError(c, err, msgNotFound)
		return
	}

	if req.Name != "" {
		part.SetName(req.Name)
	}
	if req.Description != "" {
		part.SetDescription(req.Description)
	}
	if req.Unit != "" {
		part.SetUnit(req.Unit)
	}
	if req.Price > 0 {
		part.SetPrice(req.Price)
	}
	if req.Stock > 0 {
		part.SetStock(req.Stock)
	}

	if err := h.update.Execute(c.Request.Context(), part); err != nil {
		httputil.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, commands.ToPartResponse(part))
}
