package authuc

import (
	"context"
	"errors"
	"time"

	"github.com/fiap/postech-tc1/config"
	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/ports"
	"golang.org/x/crypto/bcrypt"
)

type Login struct {
	userRepo    ports.UserRepository
	refreshRepo ports.RefreshTokenRepository
	provider    ports.TokenProvider
	cfg         *config.Config
	// dummyHash gerado com o MESMO bcrypt cost do config pra manter tempo
	// constante quando o usuario nao existe (timing attack / user enumeration).
	dummyHash []byte
}

func NewLogin(
	userRepo ports.UserRepository,
	refreshRepo ports.RefreshTokenRepository,
	provider ports.TokenProvider,
	cfg *config.Config,
) *Login {
	dummy, _ := bcrypt.GenerateFromPassword([]byte("timing-safe-dummy"), cfg.BcryptCost)
	return &Login{
		userRepo:    userRepo,
		refreshRepo: refreshRepo,
		provider:    provider,
		cfg:         cfg,
		dummyHash:   dummy,
	}
}

func (uc *Login) Execute(ctx context.Context, email, password string) (string, string, error) {
	// 1. Buscar usuario por email
	u, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, domainerrors.ErrNotFound) {
		// Erro de infraestrutura — propaga pra 500
		return "", "", err
	}

	// 2. Comparar senha com bcrypt — SEMPRE, mesmo se user nao existe.
	// Mantem tempo constante e previne user enumeration via timing.
	hashToCompare := uc.dummyHash
	if u != nil {
		hashToCompare = []byte(u.PasswordHash())
	}
	bcryptErr := bcrypt.CompareHashAndPassword(hashToCompare, []byte(password))

	// 3. User nao existe → erro generico
	if u == nil {
		return "", "", domainerrors.ErrInvalidCredentials
	}

	// 4. Conta bloqueada → erro generico (nao revelamos estado de bloqueio
	// pra prevenir enumeration/reconhaissance)
	if u.IsLocked() {
		return "", "", domainerrors.ErrInvalidCredentials
	}

	// 5. Senha errada → incrementa contador e bloqueia se atingir limite
	if bcryptErr != nil {
		lockDuration := time.Duration(uc.cfg.LoginLockMin) * time.Minute
		u.RegisterFailedLogin(uc.cfg.MaxFailedLogins, lockDuration)
		// Best-effort: se falhar o update nao trava o login, mas loga
		_ = uc.userRepo.Update(ctx, u)
		return "", "", domainerrors.ErrInvalidCredentials
	}

	// 6. Login bem-sucedido → reseta contador de falhas
	u.ResetFailedLogins()
	if err := uc.userRepo.Update(ctx, u); err != nil {
		return "", "", err
	}

	// 7. Gerar access token
	accessToken, err := uc.provider.GenerateAccessToken(u)
	if err != nil {
		return "", "", err
	}

	// 8. Gerar refresh token (raw pro cliente, hash pro banco)
	rawRefresh, hashRefresh, err := uc.provider.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	// 9. Persistir refresh token
	expiresAt := time.Now().Add(uc.provider.RefreshTokenExpiration())
	rt := entities.NewRefreshToken(u.ID(), hashRefresh, expiresAt)
	if err := uc.refreshRepo.Create(ctx, rt); err != nil {
		return "", "", err
	}

	return accessToken, rawRefresh, nil
}
