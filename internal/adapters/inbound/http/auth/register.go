package authhandler

import (
	"net/http"

	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/gin-gonic/gin"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
)

// Register godoc
// @Summary     Registrar usuario
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body commands.RegisterRequest true "Dados do usuario"
// @Success     201 {object} commands.MessageResponse
// @Router      /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req commands.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.HandleBadRequest(c, err)
		return
	}

	err := h.register.Execute(c.Request.Context(), req.Name, req.Email, req.Password, entities.RoleClient)
	if err != nil {
		httputil.HandleError(c, err, "email already registered")
		return
	}

	c.JSON(http.StatusCreated, commands.MessageResponse{Message: "user registered successfully"})
}
