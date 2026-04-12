// Package token provê a interface e implementação de geração de tokens JWT e refresh tokens.
// A interface Service permite substituir a implementação sem alterar os use cases.
package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/fiap/postech-tc1/internal/domain/user"
	"github.com/golang-jwt/jwt/v5"
)

// Service abstrai a geração de tokens de autenticação.
// Use cases dependem desta interface, mantendo-os desacoplados de JWT e crypto.
type Service interface {
	// GenerateAccessToken cria um JWT assinado para o usuario com o tempo de expiracao fornecido.
	GenerateAccessToken(u *user.User, expiresIn time.Duration) (string, error)

	// GenerateRefreshToken gera um par (raw, hash) para refresh token.
	// raw e enviado ao cliente; hash e armazenado no banco.
	GenerateRefreshToken() (raw string, hash string, err error)

	// HashToken retorna o SHA-256 hex-encoded de um token raw.
	// Usado para derivar o hash ao receber o token do cliente.
	HashToken(raw string) string
}

type jwtService struct {
	secret string
}

// New cria uma implementacao de Service usando HMAC-SHA256 para JWT
// e crypto/rand para refresh tokens.
func New(secret string) Service {
	return &jwtService{secret: secret}
}

func (s *jwtService) GenerateAccessToken(u *user.User, expiresIn time.Duration) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":   u.ID(),
		"email": u.Email(),
		"role":  string(u.Role()),
		"iat":   now.Unix(),
		"exp":   now.Add(expiresIn).Unix(),
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
