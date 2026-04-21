package jwt

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/golang-jwt/jwt/v5"
)

func makeTestUser() *entities.User {
	u := entities.NewUser("Alice", "alice@test.com", "hash", entities.RoleAdmin)
	u.SetID("user-123")
	return u
}

// --- GenerateAccessToken ---

func TestTokenService_GenerateAccessToken_Claims(t *testing.T) {
	svc := New("secret", 15, 7)
	u := makeTestUser()

	tokenStr, err := svc.GenerateAccessToken(u)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	parsed, err := jwt.Parse(tokenStr, func(*jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil || !parsed.Valid {
		t.Fatalf("expected valid token, got err=%v valid=%v", err, parsed.Valid)
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("expected jwt.MapClaims")
	}

	if claims["sub"] != "user-123" {
		t.Errorf("expected sub %q, got %q", "user-123", claims["sub"])
	}
	if claims["email"] != "alice@test.com" {
		t.Errorf("expected email %q, got %q", "alice@test.com", claims["email"])
	}
	if claims["role"] != string(entities.RoleAdmin) {
		t.Errorf("expected role %q, got %q", entities.RoleAdmin, claims["role"])
	}
}

func TestTokenService_GenerateAccessToken_Expiration(t *testing.T) {
	svc := New("secret", 15, 7)

	before := time.Now().Unix()
	tokenStr, err := svc.GenerateAccessToken(makeTestUser())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	after := time.Now().Unix()

	parsed, _ := jwt.Parse(tokenStr, func(*jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	claims := parsed.Claims.(jwt.MapClaims)

	iat := int64(claims["iat"].(float64))
	exp := int64(claims["exp"].(float64))

	if iat < before || iat > after {
		t.Errorf("iat %d out of range [%d, %d]", iat, before, after)
	}
	// exp deve ser iat + accessExpMin minutos (15 * 60 = 900s)
	if exp-iat != 15*60 {
		t.Errorf("expected exp-iat = 900s, got %d", exp-iat)
	}
}

func TestTokenService_GenerateAccessToken_WrongSecretFailsParse(t *testing.T) {
	svc := New("correct-secret", 15, 7)

	tokenStr, err := svc.GenerateAccessToken(makeTestUser())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = jwt.Parse(tokenStr, func(*jwt.Token) (interface{}, error) {
		return []byte("wrong-secret"), nil
	})
	if err == nil {
		t.Fatal("expected error parsing with wrong secret, got nil")
	}
}

// --- GenerateRefreshToken ---

func TestTokenService_GenerateRefreshToken_RawIsHex32Bytes(t *testing.T) {
	svc := New("secret", 15, 7)

	raw, _, err := svc.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// 32 bytes aleatorios em hex = 64 caracteres.
	if len(raw) != 64 {
		t.Errorf("expected raw of length 64, got %d", len(raw))
	}
	if _, err := hex.DecodeString(raw); err != nil {
		t.Errorf("expected raw to be valid hex, got %v", err)
	}
}

func TestTokenService_GenerateRefreshToken_HashMatchesRaw(t *testing.T) {
	svc := New("secret", 15, 7)

	raw, hash, err := svc.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	want := sha256.Sum256([]byte(raw))
	if hash != hex.EncodeToString(want[:]) {
		t.Errorf("hash does not match sha256(raw)")
	}
}

func TestTokenService_GenerateRefreshToken_IsRandom(t *testing.T) {
	svc := New("secret", 15, 7)

	r1, _, err := svc.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	r2, _, err := svc.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if r1 == r2 {
		t.Fatal("two calls returned the same raw token — generator is not random")
	}
}

// --- HashToken ---

func TestTokenService_HashToken_Deterministic(t *testing.T) {
	svc := New("secret", 15, 7)

	h1 := svc.HashToken("input-xyz")
	h2 := svc.HashToken("input-xyz")
	if h1 != h2 {
		t.Errorf("expected deterministic hash, got %q and %q", h1, h2)
	}
}

// --- RefreshTokenExpiration ---

func TestTokenService_RefreshTokenExpiration(t *testing.T) {
	svc := New("secret", 15, 7)

	if got := svc.RefreshTokenExpiration(); got != 7*24*time.Hour {
		t.Errorf("expected 7 days, got %v", got)
	}
}
