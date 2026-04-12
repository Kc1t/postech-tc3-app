package ports

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/entities"
)

//go:generate mockgen -source=./repositories.go -destination=./mocks/repositories.go -package=mocks

// CustomerRepository define as operacoes de persistencia para clientes.
type CustomerRepository interface {
	Create(ctx context.Context, c *entities.Customer) error
	FindByID(ctx context.Context, id string) (*entities.Customer, error)
	FindByDocument(ctx context.Context, document string) (*entities.Customer, error)
	FindAll(ctx context.Context) ([]*entities.Customer, error)
	Update(ctx context.Context, c *entities.Customer) error
	Delete(ctx context.Context, id string) error
}

// VehicleRepository define as operacoes de persistencia para veiculos.
type VehicleRepository interface {
	Create(ctx context.Context, v *entities.Vehicle) error
	FindByID(ctx context.Context, id string) (*entities.Vehicle, error)
	FindByCustomerID(ctx context.Context, customerID string) ([]*entities.Vehicle, error)
	FindAll(ctx context.Context) ([]*entities.Vehicle, error)
	Update(ctx context.Context, v *entities.Vehicle) error
	Delete(ctx context.Context, id string) error
}

// ServiceOrderRepository define as operacoes de persistencia para ordens de servico.
type ServiceOrderRepository interface {
	Create(ctx context.Context, so *entities.ServiceOrder) error
	FindByID(ctx context.Context, id string) (*entities.ServiceOrder, error)
	FindAll(ctx context.Context) ([]*entities.ServiceOrder, error)
	FindByCustomerID(ctx context.Context, customerID string) ([]*entities.ServiceOrder, error)
	UpdateStatus(ctx context.Context, id string, status entities.OrderStatus) error
	Update(ctx context.Context, so *entities.ServiceOrder) error
	Delete(ctx context.Context, id string) error
}

// ServiceRepository define as operacoes de persistencia para servicos.
type ServiceRepository interface {
	Create(ctx context.Context, s *entities.Service) error
	FindByID(ctx context.Context, id string) (*entities.Service, error)
	FindAll(ctx context.Context) ([]*entities.Service, error)
	Update(ctx context.Context, s *entities.Service) error
	Delete(ctx context.Context, id string) error
}

// PartRepository define as operacoes de persistencia para pecas e insumos.
type PartRepository interface {
	Create(ctx context.Context, p *entities.Part) error
	FindByID(ctx context.Context, id string) (*entities.Part, error)
	FindAll(ctx context.Context) ([]*entities.Part, error)
	Update(ctx context.Context, p *entities.Part) error
	Delete(ctx context.Context, id string) error
	UpdateStock(ctx context.Context, id string, delta int) error
}
