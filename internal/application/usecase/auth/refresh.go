package authuc

import (
	"context"
	"errors"
	"time"

	"github.com/fiap/postech-tc1/config"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
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
	tokenHash := hashToken(rawRefreshToken)

	var accessToken string
	var newRaw string

	// Find + validate + revoke + create dentro da mesma transacao.
	// O Revoke do repositorio usa WHERE id = ? AND revoked = false e checa
	// RowsAffected, garantindo que apenas UMA requisicao consegue consumir
	// o token em cenarios de concorrencia (single-use real).
	err := uc.txManager.RunInTransaction(ctx, func(txCtx context.Context) error {
		rt, err := uc.refreshRepo.FindByTokenHash(txCtx, tokenHash)
		if err != nil {
			if errors.Is(err, domainerrors.ErrNotFound) {
				return ErrInvalidRefreshToken
			}
			return err
		}

		if !rt.IsValid() {
			return ErrInvalidRefreshToken
		}

		// Revoke eh quem efetivamente "consome" o token.
		// Se outra requisicao concorrente ja revogou, aqui recebemos
		// ErrNotFound (RowsAffected = 0) e retornamos token invalido.
		if err := uc.refreshRepo.Revoke(txCtx, rt.ID()); err != nil {
			if errors.Is(err, domainerrors.ErrNotFound) {
				return ErrInvalidRefreshToken
			}
			return err
		}

		u, err := uc.userRepo.FindByID(txCtx, rt.UserID())
		if err != nil {
			return err
		}

		accessToken, err = generateAccessToken(u, uc.cfg.JWTSecret, uc.cfg.AccessTokenExpMin)
		if err != nil {
			return err
		}

		rawToken, hash, err := generateRefreshToken()
		if err != nil {
			return err
		}

		expiresAt := time.Now().Add(time.Duration(uc.cfg.RefreshTokenExpDays) * 24 * time.Hour)
		newRT := user.NewRefreshToken(u.ID(), hash, expiresAt)
		if err := uc.refreshRepo.Create(txCtx, newRT); err != nil {
			return err
		}

		newRaw = rawToken
		return nil
	})

	if err != nil {
		return "", "", err
	}

	return accessToken, newRaw, nil
}
