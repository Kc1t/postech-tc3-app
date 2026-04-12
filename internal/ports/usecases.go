package ports

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/customer"
	"github.com/fiap/postech-tc1/internal/domain/part"
	"github.com/fiap/postech-tc1/internal/domain/service"
	"github.com/fiap/postech-tc1/internal/domain/serviceorder"
	"github.com/fiap/postech-tc1/internal/domain/user"
	"github.com/fiap/postech-tc1/internal/domain/vehicle"
)

//go:generate mockgen -source=./usecases.go -destination=./mocks/usecases.go -package=mocks

// --- Auth ---

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

// --- Customer ---

type CreateCustomerUseCase interface {
	Execute(ctx context.Context, c *customer.Customer) error
}

type GetCustomerUseCase interface {
	Execute(ctx context.Context, id string) (*customer.Customer, error)
}

type GetCustomerByDocumentUseCase interface {
	Execute(ctx context.Context, document string) (*customer.Customer, error)
}

type ListCustomersUseCase interface {
	Execute(ctx context.Context) ([]*customer.Customer, error)
}

type UpdateCustomerUseCase interface {
	Execute(ctx context.Context, c *customer.Customer) error
}

type DeleteCustomerUseCase interface {
	Execute(ctx context.Context, id string) error
}

// --- Vehicle ---

type CreateVehicleUseCase interface {
	Execute(ctx context.Context, v *vehicle.Vehicle) error
}

type GetVehicleUseCase interface {
	Execute(ctx context.Context, id string) (*vehicle.Vehicle, error)
}

type ListVehiclesUseCase interface {
	Execute(ctx context.Context) ([]*vehicle.Vehicle, error)
}

type ListVehiclesByCustomerUseCase interface {
	Execute(ctx context.Context, customerID string) ([]*vehicle.Vehicle, error)
}

type UpdateVehicleUseCase interface {
	Execute(ctx context.Context, v *vehicle.Vehicle) error
}

type DeleteVehicleUseCase interface {
	Execute(ctx context.Context, id string) error
}

// --- ServiceOrder ---

type CreateServiceOrderUseCase interface {
	Execute(ctx context.Context, so *serviceorder.ServiceOrder) error
}

type GetServiceOrderUseCase interface {
	Execute(ctx context.Context, id string) (*serviceorder.ServiceOrder, error)
}

type ListServiceOrdersUseCase interface {
	Execute(ctx context.Context) ([]*serviceorder.ServiceOrder, error)
}

type ListServiceOrdersByCustomerUseCase interface {
	Execute(ctx context.Context, customerID string) ([]*serviceorder.ServiceOrder, error)
}

type UpdateServiceOrderStatusUseCase interface {
	Execute(ctx context.Context, id string, status serviceorder.Status) error
}

type UpdateServiceOrderUseCase interface {
	Execute(ctx context.Context, so *serviceorder.ServiceOrder) error
}

type DeleteServiceOrderUseCase interface {
	Execute(ctx context.Context, id string) error
}

// --- Service ---

type CreateServiceUseCase interface {
	Execute(ctx context.Context, s *service.Service) error
}

type GetServiceUseCase interface {
	Execute(ctx context.Context, id string) (*service.Service, error)
}

type ListServicesUseCase interface {
	Execute(ctx context.Context) ([]*service.Service, error)
}

type UpdateServiceUseCase interface {
	Execute(ctx context.Context, s *service.Service) error
}

type DeleteServiceUseCase interface {
	Execute(ctx context.Context, id string) error
}

// --- Part ---

type CreatePartUseCase interface {
	Execute(ctx context.Context, p *part.Part) error
}

type GetPartUseCase interface {
	Execute(ctx context.Context, id string) (*part.Part, error)
}

type ListPartsUseCase interface {
	Execute(ctx context.Context) ([]*part.Part, error)
}

type UpdatePartUseCase interface {
	Execute(ctx context.Context, p *part.Part) error
}

type DeletePartUseCase interface {
	Execute(ctx context.Context, id string) error
}

type AdjustPartStockUseCase interface {
	Execute(ctx context.Context, id string, delta int) error
}
