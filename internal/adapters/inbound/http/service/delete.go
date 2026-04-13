package servicehandler

import (
	"errors"
	"net/http"

	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/gin-gonic/gin"
)

// Delete godoc
// @Summary     Deletar servico
// @Tags        services
// @Param       id path string true "Service ID"
// @Success     204
// @Security    BearerAuth
// @Router      /services/{id} [delete]
func (h *ServiceHandler) Delete(c *gin.Context) {
	if err := h.delete.Execute(c.Request.Context(), c.Param("id")); err != nil {
		if errors.Is(err, domainerrors.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
