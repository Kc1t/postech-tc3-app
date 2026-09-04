package authhandler

import (
	"net/http"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	"github.com/gin-gonic/gin"
)

// Login godoc
// @Summary     Autenticar usuario
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body commands.LoginRequest true "Credenciais"
// @Success     200 {object} commands.AuthResponse
// @Router      /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req commands.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.HandleBadRequest(c, err)
		return
	}

	accessToken, refreshToken, err := h.login.Execute(c.Request.Context(), req.Email, req.Password)
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
