package authuc

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/fiap/postech-tc1/internal/domain/user"
	"github.com/golang-jwt/jwt/v5"
)

// generateAccessToken cria um JWT assinado com HMAC-SHA256.
func generateAccessToken(u *user.User, secret string, expMinutes int) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":   u.ID(),
		"email": u.Email(),
		"role":  string(u.Role()),
		"iat":   now.Unix(),
		"exp":   now.Add(time.Duration(expMinutes) * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// generateRefreshToken gera um token opaco de 32 bytes (hex-encoded).
// Retorna o token raw (enviado ao cliente) e o hash SHA-256 (salvo no banco).
func generateRefreshToken() (raw string, hash string, err error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}
	raw = hex.EncodeToString(b)
	hash = hashToken(raw)
	return raw, hash, nil
}

// hashToken retorna o SHA-256 hex-encoded de um token.
func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
