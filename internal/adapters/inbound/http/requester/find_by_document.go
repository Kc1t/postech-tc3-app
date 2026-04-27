package requesterhandler

import (
	"net/http"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	"github.com/gin-gonic/gin"
)

// FindByDocument godoc
// @Summary     Buscar cliente por documento
// @Tags        requesters
// @Produce     json
// @Param       document path string true "Requester Document"
// @Security    BearerAuth
// @Router      /requesters/document/{document} [get]
func (h *RequesterHandler) FindByDocument(c *gin.Context) {
	requester, err := h.getByDocument.Execute(c.Request.Context(), c.Param("document"))
	if err != nil {
		httputil.HandleError(c, err, msgNotFound)
		return
	}
	c.JSON(http.StatusOK, commands.ToRequesterResponse(requester))
}
