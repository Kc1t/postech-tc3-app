package authhandler

import (
	"github.com/fiap/postech-tc1/internal/ports"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	register     ports.RegisterUseCase
	login        ports.LoginUseCase
	refresh      ports.RefreshTokenUseCase
	logout       ports.LogoutUseCase
	accessExpMin int
}

func NewAuthHandler(
	register ports.RegisterUseCase,
	login ports.LoginUseCase,
	refresh ports.RefreshTokenUseCase,
	logout ports.LogoutUseCase,
	accessExpMin int,
) *AuthHandler {
	return &AuthHandler{
		register:     register,
		login:        login,
		refresh:      refresh,
		logout:       logout,
		accessExpMin: accessExpMin,
	}
}

// SetupPublicRoutes registra rotas publicas (sem JWT).
func (h *AuthHandler) SetupPublicRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
	auth.POST("/refresh", h.Refresh)
}

// SetupProtectedRoutes registra rotas que exigem JWT.
func (h *AuthHandler) SetupProtectedRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	auth.POST("/logout", h.Logout)
}
