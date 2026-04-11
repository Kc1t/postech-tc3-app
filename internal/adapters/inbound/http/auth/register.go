package authhandler

import (
	"errors"
	"net/http"

	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/domain/user"
	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary     Registrar usuario
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body commands.RegisterRequest true "Dados do usuario"
// @Success     201 {object} commands.MessageResponse
// @Failure     400 {object} commands.MessageResponse
// @Failure     409 {object} commands.MessageResponse
// @Router      /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req commands.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.register.Execute(c.Request.Context(), req.Name, req.Email, req.Password, user.RoleClient)
	if err != nil {
		if errors.Is(err, domainerrors.ErrAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, commands.MessageResponse{Message: "user registered successfully"})
}
