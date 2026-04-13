package parthandler

import (
	"net/http"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
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
		httputil.HandleError(c, err, msgNotFound)
		return
	}
	c.JSON(http.StatusOK, commands.ToPartResponse(part))
}
