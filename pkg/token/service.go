// Package token implementa a geracao de tokens JWT e refresh tokens.
// A interface TokenService fica em ports/, aqui apenas a implementacao.
package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/golang-jwt/jwt/v5"
)

type jwtService struct {
	secret         string
	accessExpMin   int
	refreshExpDays int
}

// New cria uma implementacao de TokenService usando HMAC-SHA256 para JWT
// e crypto/rand para refresh tokens.
func New(secret string, accessExpMin, refreshExpDays int) *jwtService {
	return &jwtService{
		secret:         secret,
		accessExpMin:   accessExpMin,
		refreshExpDays: refreshExpDays,
	}
}

func (s *jwtService) GenerateAccessToken(u *entities.User) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":   u.ID(),
		"email": u.Email(),
		"role":  string(u.Role()),
		"iat":   now.Unix(),
		"exp":   now.Add(time.Duration(s.accessExpMin) * time.Minute).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(s.secret))
}

func (s *jwtService) GenerateRefreshToken() (string, string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}
	raw := hex.EncodeToString(b)
	hash := s.HashToken(raw)
	return raw, hash, nil
}

func (s *jwtService) HashToken(raw string) string {
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:])
}

func (s *jwtService) RefreshTokenExpiration() time.Duration {
	return time.Duration(s.refreshExpDays) * 24 * time.Hour
}
