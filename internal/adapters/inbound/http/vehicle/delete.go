package vehiclehandler

import (
	"net/http"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
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
		httputil.HandleError(c, err, msgNotFound)
		return
	}
	c.Status(http.StatusNoContent)
}
