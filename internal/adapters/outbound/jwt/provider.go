// Package jwt implementa a geracao e validacao de tokens JWT (adapter outbound).
// Encapsula o pacote golang-jwt pra que use cases nao dependam de libs externas.
package jwt

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	jwtlib "github.com/golang-jwt/jwt/v5"
)

type Provider struct {
	secret          string
	accessExpMin    int
	refreshExpDays  int
}

func NewProvider(secret string, accessExpMin, refreshExpDays int) *Provider {
	return &Provider{
		secret:         secret,
		accessExpMin:   accessExpMin,
		refreshExpDays: refreshExpDays,
	}
}

// GenerateAccessToken cria um JWT assinado com HMAC-SHA256.
func (p *Provider) GenerateAccessToken(u *entities.User) (string, error) {
	now := time.Now()
	claims := jwtlib.MapClaims{
		"sub":   u.ID(),
		"email": u.Email(),
		"role":  string(u.Role()),
		"iat":   now.Unix(),
		"exp":   now.Add(time.Duration(p.accessExpMin) * time.Minute).Unix(),
	}
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return token.SignedString([]byte(p.secret))
}

// GenerateRefreshToken gera um token opaco de 32 bytes (hex-encoded).
// Retorna o token raw (enviado ao cliente) e o hash SHA-256 (salvo no banco).
func (p *Provider) GenerateRefreshToken() (raw string, hash string, err error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}
	raw = hex.EncodeToString(b)
	hash = HashToken(raw)
	return raw, hash, nil
}

// RefreshTokenExpiration retorna a duracao de expiracao do refresh token.
func (p *Provider) RefreshTokenExpiration() time.Duration {
	return time.Duration(p.refreshExpDays) * 24 * time.Hour
}

// HashToken retorna o SHA-256 hex-encoded de um token.
func HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
