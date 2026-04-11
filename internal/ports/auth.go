package ports

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/user"
)

// TransactionManager abstrai transacoes de banco sem expor o ORM.
// Use cases recebem essa interface pra executar multiplas operacoes
// atomicamente sem conhecer o driver/ORM subjacente.
type TransactionManager interface {
	RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// --- Repositories ---

// UserRepository define as operacoes de persistencia para usuarios.
type UserRepository interface {
	Create(ctx context.Context, u *user.User) error
	FindByEmail(ctx context.Context, email string) (*user.User, error)
	FindByID(ctx context.Context, id string) (*user.User, error)
	Update(ctx context.Context, u *user.User) error
}

// RefreshTokenRepository define as operacoes de persistencia para refresh tokens.
type RefreshTokenRepository interface {
	Create(ctx context.Context, rt *user.RefreshToken) error
	FindByTokenHash(ctx context.Context, hash string) (*user.RefreshToken, error)
	Revoke(ctx context.Context, id string) error
	RevokeByUserID(ctx context.Context, userID string) error
}

// --- Use Cases ---

// RegisterUseCase cria um novo usuario no sistema.
type RegisterUseCase interface {
	Execute(ctx context.Context, name, email, rawPassword string, role user.Role) error
}

// LoginUseCase autentica um usuario e retorna access + refresh tokens.
type LoginUseCase interface {
	Execute(ctx context.Context, email, password string) (accessToken, refreshToken string, err error)
}

// RefreshTokenUseCase rotaciona o par de tokens (refresh token rotation).
type RefreshTokenUseCase interface {
	Execute(ctx context.Context, rawRefreshToken string) (accessToken, newRefreshToken string, err error)
}

// LogoutUseCase invalida todos os refresh tokens do usuario.
type LogoutUseCase interface {
	Execute(ctx context.Context, userID string) error
}
