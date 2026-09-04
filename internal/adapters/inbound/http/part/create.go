package parthandler

import (
	"net/http"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	"github.com/gin-gonic/gin"
)

// Create godoc
// @Summary     Cadastrar peca/insumo
// @Tags        parts
// @Accept      json
// @Produce     json
// @Param       body body commands.CreatePartRequest true "Dados da peca"
// @Success     201 {object} commands.PartResponse
// @Security    BearerAuth
// @Router      /parts [post]
func (h *PartHandler) Create(c *gin.Context) {
	var req commands.CreatePartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.HandleBadRequest(c, err)
		return
	}

	part := req.ToDomain()
	if err := h.create.Execute(c.Request.Context(), part); err != nil {
		httputil.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, commands.ToPartResponse(part))
}
