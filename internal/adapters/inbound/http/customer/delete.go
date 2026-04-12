package customerhandler

import (
	"errors"
	"net/http"

	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/gin-gonic/gin"
)

// Delete godoc
// @Summary     Deletar cliente
// @Tags        customers
// @Param       id path string true "Customer ID"
// @Success     204
// @Security    BearerAuth
// @Router      /customers/{id} [delete]
func (h *CustomerHandler) Delete(c *gin.Context) {
	if err := h.delete.Execute(c.Request.Context(), c.Param("id")); err != nil {
		if errors.Is(err, domainerrors.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
