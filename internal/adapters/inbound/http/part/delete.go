package parthandler

import (
	"errors"
	"net/http"

	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/gin-gonic/gin"
)

// Delete godoc
// @Summary     Deletar peca/insumo
// @Tags        parts
// @Param       id path string true "Part ID"
// @Success     204
// @Security    BearerAuth
// @Router      /parts/{id} [delete]
func (h *PartHandler) Delete(c *gin.Context) {
	if err := h.delete.Execute(c.Request.Context(), c.Param("id")); err != nil {
		if errors.Is(err, domainerrors.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "part not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
