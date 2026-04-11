package authhandler

import (
	"errors"
	"net/http"

	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	authuc "github.com/fiap/postech-tc1/internal/application/usecase/auth"
	"github.com/gin-gonic/gin"
)

// Refresh godoc
// @Summary     Renovar tokens (refresh token rotation)
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body commands.RefreshRequest true "Refresh token"
// @Success     200 {object} commands.AuthResponse
// @Failure     400 {object} commands.MessageResponse
// @Failure     401 {object} commands.MessageResponse
// @Router      /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req commands.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := h.refresh.Execute(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, authuc.ErrInvalidRefreshToken) {
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
