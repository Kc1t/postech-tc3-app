package authuc

import (
	"context"
	"errors"
	"time"

	"github.com/fiap/postech-tc1/config"
	"github.com/fiap/postech-tc1/internal/domain/user"
	"github.com/fiap/postech-tc1/internal/ports"
)

var (
	ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
)

type Refresh struct {
	userRepo    ports.UserRepository
	refreshRepo ports.RefreshTokenRepository
	txManager   ports.TransactionManager
	cfg         *config.Config
}

func NewRefresh(
	userRepo ports.UserRepository,
	refreshRepo ports.RefreshTokenRepository,
	txManager ports.TransactionManager,
	cfg *config.Config,
) *Refresh {
	return &Refresh{userRepo: userRepo, refreshRepo: refreshRepo, txManager: txManager, cfg: cfg}
}

func (uc *Refresh) Execute(ctx context.Context, rawRefreshToken string) (string, string, error) {
	// 1. Hash do token recebido pra buscar no banco
	tokenHash := hashToken(rawRefreshToken)

	// 2. Buscar refresh token pelo hash
	rt, err := uc.refreshRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return "", "", ErrInvalidRefreshToken
	}

	// 3. Validar: nao revogado e nao expirado
	if !rt.IsValid() {
		return "", "", ErrInvalidRefreshToken
	}

	// 4. Buscar usuario pra gerar novo access token
	u, err := uc.userRepo.FindByID(ctx, rt.UserID())
	if err != nil {
		return "", "", err
	}

	// 5. Gerar novos tokens
	accessToken, err := generateAccessToken(u, uc.cfg.JWTSecret, uc.cfg.AccessTokenExpMin)
	if err != nil {
		return "", "", err
	}

	newRaw, newHash, err := generateRefreshToken()
	if err != nil {
		return "", "", err
	}

	// 6. Revogar o antigo e criar o novo em transacao
	expiresAt := time.Now().Add(time.Duration(uc.cfg.RefreshTokenExpDays) * 24 * time.Hour)
	newRT := user.NewRefreshToken(u.ID(), newHash, expiresAt)

	err = uc.txManager.RunInTransaction(ctx, func(txCtx context.Context) error {
		if err := uc.refreshRepo.Revoke(txCtx, rt.ID()); err != nil {
			return err
		}
		return uc.refreshRepo.Create(txCtx, newRT)
	})
	if err != nil {
		return "", "", err
	}

	return accessToken, newRaw, nil
}
