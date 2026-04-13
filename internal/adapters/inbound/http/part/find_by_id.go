package parthandler

import (
	"errors"
	"net/http"

	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/gin-gonic/gin"
)

// FindByID godoc
// @Summary     Buscar peca/insumo por ID
// @Tags        parts
// @Produce     json
// @Param       id path string true "Part ID"
// @Success     200 {object} commands.PartResponse
// @Security    BearerAuth
// @Router      /parts/{id} [get]
func (h *PartHandler) FindByID(c *gin.Context) {
	part, err := h.getByID.Execute(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, domainerrors.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "part not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, commands.ToPartResponse(part))
}
