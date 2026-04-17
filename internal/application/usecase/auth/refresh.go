package authuc

import (
	"context"
	"errors"
	"time"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/ports"
)

type Refresh struct {
	userRepo    ports.UserRepository
	refreshRepo ports.RefreshTokenRepository
	tokenSvc    ports.TokenService
}

func NewRefresh(
	userRepo ports.UserRepository,
	refreshRepo ports.RefreshTokenRepository,
	tokenSvc ports.TokenService,
) *Refresh {
	return &Refresh{userRepo: userRepo, refreshRepo: refreshRepo, tokenSvc: tokenSvc}
}

func (uc *Refresh) Execute(ctx context.Context, rawRefreshToken string) (string, string, error) {
	tokenHash := uc.tokenSvc.HashToken(rawRefreshToken)

	// 1. Buscar o refresh token pelo hash
	rt, err := uc.refreshRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, domainerrors.ErrNotFound) {
			return "", "", domainerrors.ErrInvalidRefreshToken
		}
		return "", "", err
	}

	// 2. Validar expiracao e estado de revogacao via dominio
	if !rt.IsValid() {
		return "", "", domainerrors.ErrInvalidRefreshToken
	}

	// 3. Buscar o usuario dono do token
	u, err := uc.userRepo.FindByID(ctx, rt.UserID())
	if err != nil {
		return "", "", err
	}

	// 4. Gerar novo access token
	accessToken, err := uc.tokenSvc.GenerateAccessToken(u)
	if err != nil {
		return "", "", err
	}

	// 5. Gerar novo refresh token
	rawToken, hash, err := uc.tokenSvc.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	expiresAt := time.Now().Add(uc.tokenSvc.RefreshTokenExpiration())
	newRT := entities.NewRefreshToken(u.ID(), hash, expiresAt)

	// 6. Revogar o token antigo e persistir o novo atomicamente
	if err := uc.refreshRepo.RotateToken(ctx, rt.ID(), newRT); err != nil {
		if errors.Is(err, domainerrors.ErrNotFound) {
			return "", "", domainerrors.ErrInvalidRefreshToken
		}
		return "", "", err
	}

	return accessToken, rawToken, nil
}
