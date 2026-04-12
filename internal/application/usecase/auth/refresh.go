package authuc

import (
	"context"
	"errors"
	"time"

	"github.com/fiap/postech-tc1/config"
	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/ports"
	"github.com/fiap/postech-tc1/pkg/token"
)

type Refresh struct {
	userRepo    ports.UserRepository
	refreshRepo ports.RefreshTokenRepository
	tokenSvc    token.Service
	cfg         *config.Config
}

func NewRefresh(
	userRepo ports.UserRepository,
	refreshRepo ports.RefreshTokenRepository,
	tokenSvc token.Service,
	cfg *config.Config,
) *Refresh {
	return &Refresh{userRepo: userRepo, refreshRepo: refreshRepo, tokenSvc: tokenSvc, cfg: cfg}
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
	accessToken, err := uc.tokenSvc.GenerateAccessToken(u, time.Duration(uc.cfg.AccessTokenExpMin)*time.Minute)
	if err != nil {
		return "", "", err
	}

	// 5. Gerar novo refresh token
	rawToken, hash, err := uc.tokenSvc.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	expiresAt := time.Now().Add(time.Duration(uc.cfg.RefreshTokenExpDays) * 24 * time.Hour)
	newRT := entities.NewRefreshToken(u.ID(), hash, expiresAt)

	// 6. Revogar o token antigo e persistir o novo atomicamente.
	// O repositorio garante single-use real com WHERE revoked = false e
	// verifica RowsAffected, prevenindo double-spend em chamadas concorrentes.
	_, err = uc.refreshRepo.RotateToken(ctx, rt.ID(), newRT)
	if err != nil {
		if errors.Is(err, domainerrors.ErrNotFound) {
			return "", "", domainerrors.ErrInvalidRefreshToken
		}
		return "", "", err
	}

	return accessToken, rawToken, nil
}
