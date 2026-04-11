package authuc

import (
	"context"
	"errors"
	"time"

	"github.com/fiap/postech-tc1/config"
	"github.com/fiap/postech-tc1/internal/domain/user"
	"github.com/fiap/postech-tc1/internal/ports"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type Login struct {
	userRepo     ports.UserRepository
	refreshRepo  ports.RefreshTokenRepository
	cfg          *config.Config
}

func NewLogin(userRepo ports.UserRepository, refreshRepo ports.RefreshTokenRepository, cfg *config.Config) *Login {
	return &Login{userRepo: userRepo, refreshRepo: refreshRepo, cfg: cfg}
}

// dummyHash e usado pra manter tempo constante quando o usuario nao existe.
// Evita timing attack que permite enumerar emails validos.
var dummyHash, _ = bcrypt.GenerateFromPassword([]byte("timing-safe-dummy"), 12)

func (uc *Login) Execute(ctx context.Context, email, password string) (string, string, error) {
	// 1. Buscar usuario por email
	u, err := uc.userRepo.FindByEmail(ctx, email)

	// 2. Comparar senha com bcrypt (SEMPRE, mesmo se user nao existe)
	// Isso garante tempo constante e previne timing attack / user enumeration
	hashToCompare := string(dummyHash)
	if err == nil {
		hashToCompare = u.PasswordHash()
	}

	bcryptErr := bcrypt.CompareHashAndPassword([]byte(hashToCompare), []byte(password))

	if err != nil || bcryptErr != nil {
		return "", "", ErrInvalidCredentials
	}

	// 3. Gerar access token
	accessToken, err := generateAccessToken(u, uc.cfg.JWTSecret, uc.cfg.AccessTokenExpMin)
	if err != nil {
		return "", "", err
	}

	// 4. Gerar refresh token (raw pro cliente, hash pro banco)
	rawRefresh, hashRefresh, err := generateRefreshToken()
	if err != nil {
		return "", "", err
	}

	// 5. Persistir refresh token
	expiresAt := time.Now().Add(time.Duration(uc.cfg.RefreshTokenExpDays) * 24 * time.Hour)
	rt := user.NewRefreshToken(u.ID(), hashRefresh, expiresAt)
	if err := uc.refreshRepo.Create(ctx, rt); err != nil {
		return "", "", err
	}

	return accessToken, rawRefresh, nil
}
