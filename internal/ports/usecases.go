package ports

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/entities"
)

//go:generate mockgen -source=./usecases.go -destination=./mocks/usecases.go -package=mocks

// --- Customer ---

type CreateCustomerUseCase interface {
	Execute(ctx context.Context, c *entities.Customer) error
}

type GetCustomerUseCase interface {
	Execute(ctx context.Context, id string) (*entities.Customer, error)
}

type GetCustomerByDocumentUseCase interface {
	Execute(ctx context.Context, document string) (*entities.Customer, error)
}

type ListCustomersUseCase interface {
	Execute(ctx context.Context) ([]*entities.Customer, error)
}

type UpdateCustomerUseCase interface {
	Execute(ctx context.Context, c *entities.Customer) error
}

type DeleteCustomerUseCase interface {
	Execute(ctx context.Context, id string) error
}

// --- Vehicle ---

type CreateVehicleUseCase interface {
	Execute(ctx context.Context, v *entities.Vehicle) error
}

type GetVehicleUseCase interface {
	Execute(ctx context.Context, id string) (*entities.Vehicle, error)
}

type ListVehiclesUseCase interface {
	Execute(ctx context.Context) ([]*entities.Vehicle, error)
}

type ListVehiclesByCustomerUseCase interface {
	Execute(ctx context.Context, customerID string) ([]*entities.Vehicle, error)
}

type UpdateVehicleUseCase interface {
	Execute(ctx context.Context, v *entities.Vehicle) error
}

type DeleteVehicleUseCase interface {
	Execute(ctx context.Context, id string) error
}

// --- ServiceOrder ---

type CreateServiceOrderInput struct {
	CustomerCPF  string
	VehiclePlate string
	Notes        string
}

type CreateServiceOrderUseCase interface {
	Execute(ctx context.Context, input CreateServiceOrderInput) (*entities.ServiceOrder, error)
}

type GetServiceOrderUseCase interface {
	Execute(ctx context.Context, id string) (*entities.ServiceOrder, error)
}

type ListServiceOrdersUseCase interface {
	Execute(ctx context.Context) ([]*entities.ServiceOrder, error)
}

type ListServiceOrdersByCustomerUseCase interface {
	Execute(ctx context.Context, customerID string) ([]*entities.ServiceOrder, error)
}

type UpdateStatusInput struct {
	ID           string
	Status       entities.OrderStatus
	ServiceCodes []int                    // obrigatorio quando status = awaiting_approval
	Parts        []UpdateStatusPartInput  // obrigatorio quando status = awaiting_approval
}

type UpdateStatusPartInput struct {
	ManufacturerCode string
	Quantity         int
}

type UpdateServiceOrderStatusUseCase interface {
	Execute(ctx context.Context, input UpdateStatusInput) error
}

type GetServiceOrderByCodeUseCase interface {
	Execute(ctx context.Context, code int, customerCPF string) (*entities.ServiceOrder, error)
}

type ListServiceOrdersByCPFUseCase interface {
	Execute(ctx context.Context, cpf string) ([]*entities.ServiceOrder, error)
}

type ApproveServiceOrderUseCase interface {
	Execute(ctx context.Context, code int, customerCPF string) error
}

type RejectServiceOrderUseCase interface {
	Execute(ctx context.Context, code int, customerCPF string) error
}

type UpdateServiceOrderUseCase interface {
	Execute(ctx context.Context, so *entities.ServiceOrder) error
}

type DeleteServiceOrderUseCase interface {
	Execute(ctx context.Context, id string) error
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
