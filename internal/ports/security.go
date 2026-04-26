package ports

import (
	"time"

	"github.com/fiap/postech-tc1/internal/domain/entities"
)

// TokenService abstrai a geracao e validacao de tokens de autenticacao.
// A implementacao concreta fica em internal/adapters/outbound/jwt,
// desacoplando use cases de JWT e crypto.
type TokenService interface {
	GenerateAccessToken(u *entities.User) (string, error)
	GenerateRefreshToken() (raw string, hash string, err error)
	HashToken(raw string) string
	RefreshTokenExpiration() time.Duration
}

// PasswordHasher abstrai o hashing de senhas.
// A implementacao concreta fica em pkg/hasher, desacoplando use cases de bcrypt.
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) error
}
