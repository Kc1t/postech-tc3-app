package authuc

import (
	"context"
	"errors"
	"time"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/ports"
)

type Login struct {
	userRepo        ports.UserRepository
	refreshRepo     ports.RefreshTokenRepository
	tokenSvc        ports.TokenService
	hasher          ports.PasswordHasher
	maxFailedLogins int
	lockDurationMin int
	// dummyHash gerado com o MESMO algoritmo e custo pra manter tempo
	// constante quando o usuario nao existe (timing attack / user enumeration).
	dummyHash string
}

func NewLogin(
	userRepo ports.UserRepository,
	refreshRepo ports.RefreshTokenRepository,
	tokenSvc ports.TokenService,
	hasher ports.PasswordHasher,
	maxFailedLogins int,
	lockDurationMin int,
) *Login {
	dummy, _ := hasher.Hash("timing-safe-dummy")
	return &Login{
		userRepo:        userRepo,
		refreshRepo:     refreshRepo,
		tokenSvc:        tokenSvc,
		hasher:          hasher,
		maxFailedLogins: maxFailedLogins,
		lockDurationMin: lockDurationMin,
		dummyHash:       dummy,
	}
}

func (uc *Login) Execute(ctx context.Context, email, password string) (string, string, error) {
	// 1. Buscar usuario por email
	u, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, domainerrors.ErrNotFound) {
		return "", "", err
	}

	// 2. Comparar senha — SEMPRE, mesmo se user nao existe.
	// Mantem tempo constante e previne user enumeration via timing.
	hashToCompare := uc.dummyHash
	if u != nil {
		hashToCompare = u.PasswordHash()
	}
	compareErr := uc.hasher.Compare(hashToCompare, password)

	// 3. User nao existe → erro generico
	if u == nil {
		return "", "", domainerrors.ErrInvalidCredentials
	}

	// 4. Conta bloqueada → erro generico (nao revelamos estado de bloqueio)
	if u.IsLocked() {
		return "", "", domainerrors.ErrInvalidCredentials
	}

	// 5. Senha errada → incrementa contador e bloqueia se atingir limite
	if compareErr != nil {
		lockDuration := time.Duration(uc.lockDurationMin) * time.Minute
		u.RegisterFailedLogin(uc.maxFailedLogins, lockDuration)
		_ = uc.userRepo.Update(ctx, u)
		return "", "", domainerrors.ErrInvalidCredentials
	}

	// 6. Login bem-sucedido → reseta contador de falhas
	u.ResetFailedLogins()
	if err := uc.userRepo.Update(ctx, u); err != nil {
		return "", "", err
	}

	// 7. Gerar access token
	accessToken, err := uc.tokenSvc.GenerateAccessToken(u)
	if err != nil {
		return "", "", err
	}

	// 8. Gerar refresh token (raw pro cliente, hash pro banco)
	rawRefresh, hashRefresh, err := uc.tokenSvc.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	// 9. Persistir refresh token
	expiresAt := time.Now().Add(uc.tokenSvc.RefreshTokenExpiration())
	rt := entities.NewRefreshToken(u.ID(), hashRefresh, expiresAt)
	if err := uc.refreshRepo.Create(ctx, rt); err != nil {
		return "", "", err
	}

	return accessToken, rawRefresh, nil
}
