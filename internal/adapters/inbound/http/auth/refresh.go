package authhandler

import (
	"net/http"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	"github.com/gin-gonic/gin"
)

// Refresh godoc
// @Summary     Renovar tokens (refresh token rotation)
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body commands.RefreshRequest true "Refresh token"
// @Success     200 {object} commands.AuthResponse
// @Router      /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req commands.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.HandleBadRequest(c, err)
		return
	}

	accessToken, refreshToken, err := h.refresh.Execute(c.Request.Context(), req.RefreshToken)
	if err != nil {
		httputil.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, commands.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    h.accessExpMin * 60,
		TokenType:    "Bearer",
	})
}
