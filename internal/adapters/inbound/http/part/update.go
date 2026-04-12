package parthandler

import (
	"errors"
	"net/http"

	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	part, err := h.getByID.Execute(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, domainerrors.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "part not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, commands.ToPartResponse(part))
}
