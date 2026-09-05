package entities

import "time"

type RefreshToken struct {
	id        string
	userID    string
	tokenHash string // SHA-256 do token raw — nunca armazenar o token em texto
	expiresAt time.Time
	revoked   bool
	createdAt time.Time
}

// NewRefreshToken cria um refresh token sem ID — o UUID e gerado pelo banco
// (default:gen_random_uuid()) e o repository preenche via SetID apos Create.
// Mantem a geracao de ID na camada de infraestrutura, consistente com User e Requester.
func NewRefreshToken(userID, tokenHash string, expiresAt time.Time) *RefreshToken {
	return &RefreshToken{
		userID:    userID,
		tokenHash: tokenHash,
		expiresAt: expiresAt,
		revoked:   false,
		createdAt: time.Now(),
	}
}

// ReconstituteRefreshToken restaura a partir de dados persistidos.
func ReconstituteRefreshToken(id, userID, tokenHash string, expiresAt time.Time, revoked bool, createdAt time.Time) *RefreshToken {
	return &RefreshToken{
		id:        id,
		userID:    userID,
		tokenHash: tokenHash,
		expiresAt: expiresAt,
		revoked:   revoked,
		createdAt: createdAt,
	}
}

// Getters
func (rt *RefreshToken) ID() string           { return rt.id }
func (rt *RefreshToken) UserID() string       { return rt.userID }
func (rt *RefreshToken) TokenHash() string    { return rt.tokenHash }
func (rt *RefreshToken) ExpiresAt() time.Time { return rt.expiresAt }
func (rt *RefreshToken) Revoked() bool        { return rt.revoked }
func (rt *RefreshToken) CreatedAt() time.Time { return rt.createdAt }

// Setters
func (rt *RefreshToken) SetID(id string) { rt.id = id }

// Comportamentos
func (rt *RefreshToken) Revoke()         { rt.revoked = true }
func (rt *RefreshToken) IsExpired() bool { return time.Now().After(rt.expiresAt) }
func (rt *RefreshToken) IsValid() bool   { return !rt.revoked && !rt.IsExpired() }
