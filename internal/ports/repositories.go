package ports

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/customer"
	"github.com/fiap/postech-tc1/internal/domain/part"
	"github.com/fiap/postech-tc1/internal/domain/service"
	"github.com/fiap/postech-tc1/internal/domain/serviceorder"
	"github.com/fiap/postech-tc1/internal/domain/vehicle"
)

//go:generate mockgen -source=./repositories.go -destination=./mocks/repositories.go -package=mocks

// CustomerRepository define as operacoes de persistencia para clientes.
type CustomerRepository interface {
	Create(ctx context.Context, c *customer.Customer) error
	FindByID(ctx context.Context, id string) (*customer.Customer, error)
	FindByDocument(ctx context.Context, document string) (*customer.Customer, error)
	FindAll(ctx context.Context) ([]*customer.Customer, error)
	Update(ctx context.Context, c *customer.Customer) error
	Delete(ctx context.Context, id string) error
}

// VehicleRepository define as operacoes de persistencia para veiculos.
type VehicleRepository interface {
	Create(ctx context.Context, v *vehicle.Vehicle) error
	FindByID(ctx context.Context, id string) (*vehicle.Vehicle, error)
	FindByCustomerID(ctx context.Context, customerID string) ([]*vehicle.Vehicle, error)
	FindAll(ctx context.Context) ([]*vehicle.Vehicle, error)
	Update(ctx context.Context, v *vehicle.Vehicle) error
	Delete(ctx context.Context, id string) error
}

// ServiceOrderRepository define as operacoes de persistencia para ordens de servico.
type ServiceOrderRepository interface {
	Create(ctx context.Context, so *serviceorder.ServiceOrder) error
	FindByID(ctx context.Context, id string) (*serviceorder.ServiceOrder, error)
	FindAll(ctx context.Context) ([]*serviceorder.ServiceOrder, error)
	FindByCustomerID(ctx context.Context, customerID string) ([]*serviceorder.ServiceOrder, error)
	UpdateStatus(ctx context.Context, id string, status serviceorder.Status) error
	Update(ctx context.Context, so *serviceorder.ServiceOrder) error
	Delete(ctx context.Context, id string) error
}

// ServiceRepository define as operacoes de persistencia para servicos.
type ServiceRepository interface {
	Create(ctx context.Context, s *service.Service) error
	FindByID(ctx context.Context, id string) (*service.Service, error)
	FindAll(ctx context.Context) ([]*service.Service, error)
	Update(ctx context.Context, s *service.Service) error
	Delete(ctx context.Context, id string) error
}

// PartRepository define as operacoes de persistencia para pecas e insumos.
type PartRepository interface {
	Create(ctx context.Context, p *part.Part) error
	FindByID(ctx context.Context, id string) (*part.Part, error)
	FindAll(ctx context.Context) ([]*part.Part, error)
	Update(ctx context.Context, p *part.Part) error
	Delete(ctx context.Context, id string) error
	UpdateStock(ctx context.Context, id string, delta int) error
}
