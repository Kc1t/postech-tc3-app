package ports

import (
	"context"
	"time"

	"github.com/fiap/postech-tc1/internal/domain/entities"
)

//go:generate mockgen -source=./usecases.go -destination=./mocks/usecases.go -package=mocks

// --- Requester ---

type CreateRequesterUseCase interface {
	Execute(ctx context.Context, c *entities.Requester) error
}

type GetRequesterUseCase interface {
	Execute(ctx context.Context, id string) (*entities.Requester, error)
}

type GetRequesterByDocumentUseCase interface {
	Execute(ctx context.Context, document string) (*entities.Requester, error)
}

type ListRequestersUseCase interface {
	Execute(ctx context.Context) ([]*entities.Requester, error)
}

type UpdateRequesterUseCase interface {
	Execute(ctx context.Context, c *entities.Requester) error
}

type DeleteRequesterUseCase interface {
	Execute(ctx context.Context, id string) error
}

// --- Vehicle ---

type CreateVehicleUseCase interface {
	Execute(ctx context.Context, input entities.VehicleInput) (*entities.Vehicle, error)
}

type GetVehicleUseCase interface {
	Execute(ctx context.Context, id string) (*entities.Vehicle, error)
}

type ListVehiclesUseCase interface {
	Execute(ctx context.Context) ([]*entities.Vehicle, error)
}

type ListVehiclesByRequesterUseCase interface {
	Execute(ctx context.Context, requesterID string) ([]*entities.Vehicle, error)
}

type UpdateVehicleUseCase interface {
	Execute(ctx context.Context, v *entities.Vehicle) error
}

type DeleteVehicleUseCase interface {
	Execute(ctx context.Context, id string) error
}

// --- ServiceOrder ---

type CreateServiceOrderUseCase interface {
	Execute(ctx context.Context, input entities.ServiceOrderInput) (*entities.ServiceOrder, error)
}

type GetServiceOrderUseCase interface {
	Execute(ctx context.Context, id string) (*entities.ServiceOrder, error)
}

type ListServiceOrdersUseCase interface {
	Execute(ctx context.Context) ([]*entities.ServiceOrder, error)
}

type ListServiceOrdersByRequesterUseCase interface {
	Execute(ctx context.Context, requesterID string) ([]*entities.ServiceOrder, error)
}

type UpdateServiceOrderStatusUseCase interface {
	Execute(ctx context.Context, input entities.StatusUpdate) error
}

type GetServiceOrderByCodeUseCase interface {
	Execute(ctx context.Context, code int, requesterDocument string) (*entities.ServiceOrder, error)
}

type ListServiceOrdersByDocumentUseCase interface {
	Execute(ctx context.Context, document string) ([]*entities.ServiceOrder, error)
}

type UpdateServiceOrderStatusByCodeUseCase interface {
	Execute(ctx context.Context, code int, requesterDocument string, newStatus entities.OrderStatus) error
}

type UpdateServiceOrderUseCase interface {
	Execute(ctx context.Context, so *entities.ServiceOrder) error
}

type DeleteServiceOrderUseCase interface {
	Execute(ctx context.Context, id string) error
}

type GetAverageExecutionTimeUseCase interface {
	Execute(ctx context.Context) (time.Duration, error)
}

// --- Service ---

type CreateServiceUseCase interface {
	Execute(ctx context.Context, s *entities.Service) error
}

type GetServiceUseCase interface {
	Execute(ctx context.Context, id string) (*entities.Service, error)
}

type ListServicesUseCase interface {
	Execute(ctx context.Context) ([]*entities.Service, error)
}

type UpdateServiceUseCase interface {
	Execute(ctx context.Context, s *entities.Service) error
}

type DeleteServiceUseCase interface {
	Execute(ctx context.Context, id string) error
}

// --- Part ---

type CreatePartUseCase interface {
	Execute(ctx context.Context, p *entities.Part) error
}

type GetPartUseCase interface {
	Execute(ctx context.Context, id string) (*entities.Part, error)
}

type ListPartsUseCase interface {
	Execute(ctx context.Context) ([]*entities.Part, error)
}

type UpdatePartUseCase interface {
	Execute(ctx context.Context, p *entities.Part) error
}

type DeletePartUseCase interface {
	Execute(ctx context.Context, id string) error
}

type AdjustPartStockUseCase interface {
	Execute(ctx context.Context, id string, delta int) error
}

// --- Auth ---

// RegisterUseCase cria um novo usuario no sistema.
type RegisterUseCase interface {
	Execute(ctx context.Context, name, email, rawPassword string, role entities.Role) error
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
