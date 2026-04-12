package authhandler

import (
	"errors"
	"net/http"

	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/gin-gonic/gin"
)

// Login godoc
// @Summary     Autenticar usuario
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body commands.LoginRequest true "Credenciais"
// @Success     200 {object} commands.AuthResponse
// @Failure     400 {object} commands.MessageResponse
// @Failure     401 {object} commands.MessageResponse
// @Router      /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req commands.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := h.login.Execute(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domainerrors.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, commands.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    h.accessExpMin * 60,
		TokenType:    "Bearer",
	})
}
