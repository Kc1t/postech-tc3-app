package ports

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/customer"
	"github.com/fiap/postech-tc1/internal/domain/part"
	"github.com/fiap/postech-tc1/internal/domain/service"
	"github.com/fiap/postech-tc1/internal/domain/serviceorder"
	"github.com/fiap/postech-tc1/internal/domain/vehicle"
)

// CustomerUseCase define os casos de uso para clientes.
type CustomerUseCase interface {
	Create(ctx context.Context, c *customer.Customer) error
	GetByID(ctx context.Context, id string) (*customer.Customer, error)
	GetByDocument(ctx context.Context, document string) (*customer.Customer, error)
	GetAll(ctx context.Context) ([]*customer.Customer, error)
	Update(ctx context.Context, c *customer.Customer) error
	Delete(ctx context.Context, id string) error
}

// VehicleUseCase define os casos de uso para veiculos.
type VehicleUseCase interface {
	Create(ctx context.Context, v *vehicle.Vehicle) error
	GetByID(ctx context.Context, id string) (*vehicle.Vehicle, error)
	GetByCustomerID(ctx context.Context, customerID string) ([]*vehicle.Vehicle, error)
	GetAll(ctx context.Context) ([]*vehicle.Vehicle, error)
	Update(ctx context.Context, v *vehicle.Vehicle) error
	Delete(ctx context.Context, id string) error
}

// ServiceOrderUseCase define os casos de uso para ordens de servico.
type ServiceOrderUseCase interface {
	Create(ctx context.Context, so *serviceorder.ServiceOrder) error
	GetByID(ctx context.Context, id string) (*serviceorder.ServiceOrder, error)
	GetAll(ctx context.Context) ([]*serviceorder.ServiceOrder, error)
	GetByCustomerID(ctx context.Context, customerID string) ([]*serviceorder.ServiceOrder, error)
	UpdateStatus(ctx context.Context, id string, status serviceorder.Status) error
	Update(ctx context.Context, so *serviceorder.ServiceOrder) error
	Delete(ctx context.Context, id string) error
}

// ServiceUseCase define os casos de uso para servicos.
type ServiceUseCase interface {
	Create(ctx context.Context, s *service.Service) error
	GetByID(ctx context.Context, id string) (*service.Service, error)
	GetAll(ctx context.Context) ([]*service.Service, error)
	Update(ctx context.Context, s *service.Service) error
	Delete(ctx context.Context, id string) error
}

// PartUseCase define os casos de uso para pecas e insumos.
type PartUseCase interface {
	Create(ctx context.Context, p *part.Part) error
	GetByID(ctx context.Context, id string) (*part.Part, error)
	GetAll(ctx context.Context) ([]*part.Part, error)
	Update(ctx context.Context, p *part.Part) error
	Delete(ctx context.Context, id string) error
	AdjustStock(ctx context.Context, id string, delta int) error
}
