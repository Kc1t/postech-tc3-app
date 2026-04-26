package authhandler

import (
	"github.com/fiap/postech-tc1/internal/ports"
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
