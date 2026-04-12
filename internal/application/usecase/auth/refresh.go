package authuc

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/ports"
)

type Refresh struct {
	userRepo    ports.UserRepository
	refreshRepo ports.RefreshTokenRepository
	provider    ports.TokenProvider
}

func NewRefresh(
	userRepo ports.UserRepository,
	refreshRepo ports.RefreshTokenRepository,
	provider ports.TokenProvider,
) *Refresh {
	return &Refresh{userRepo: userRepo, refreshRepo: refreshRepo, provider: provider}
}

func (uc *Refresh) Execute(ctx context.Context, rawRefreshToken string) (string, string, error) {
	tokenHash := hashToken(rawRefreshToken)

	// 1. Gerar novo refresh token
	newRaw, newHash, err := uc.provider.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	expiresAt := time.Now().Add(uc.provider.RefreshTokenExpiration())
	newRT := entities.NewRefreshToken("", newHash, expiresAt) // userID preenchido pelo repo

	// 2. Rotacionar atomicamente (find + validate + revoke + create)
	// O RotateToken precisa do userID do token antigo, entao fazemos em duas etapas:
	// primeiro buscamos o token antigo pra pegar o userID.
	oldRT, err := uc.refreshRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		if isNotFound(err) {
			return "", "", domainerrors.ErrInvalidRefreshToken
		}
		return "", "", err
	}

	if !oldRT.IsValid() {
		return "", "", domainerrors.ErrInvalidRefreshToken
	}

	// Criar o novo RT com o userID correto
	newRT = entities.NewRefreshToken(oldRT.UserID(), newHash, expiresAt)

	_, err = uc.refreshRepo.RotateToken(ctx, tokenHash, newRT)
	if err != nil {
		return "", "", err
	}

	// 3. Gerar access token
	u, err := uc.userRepo.FindByID(ctx, oldRT.UserID())
	if err != nil {
		return "", "", err
	}

	accessToken, err := uc.provider.GenerateAccessToken(u)
	if err != nil {
		return "", "", err
	}

	return accessToken, newRaw, nil
}

func isNotFound(err error) bool {
	return err == domainerrors.ErrNotFound
}

// hashToken retorna o SHA-256 hex-encoded de um token.
func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
