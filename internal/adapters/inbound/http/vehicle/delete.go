package vehiclehandler

import (
	"errors"
	"net/http"

	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/gin-gonic/gin"
)

// Delete godoc
// @Summary     Deletar veiculo
// @Tags        vehicles
// @Param       id path string true "Vehicle ID"
// @Success     204
// @Security    BearerAuth
// @Router      /vehicles/{id} [delete]
func (h *VehicleHandler) Delete(c *gin.Context) {
	if err := h.delete.Execute(c.Request.Context(), c.Param("id")); err != nil {
		if errors.Is(err, domainerrors.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "vehicle not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
