package serviceorderuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/ports"
)

type CreateServiceOrder struct {
	repo         ports.ServiceOrderRepository
	customerRepo ports.CustomerRepository
	vehicleRepo  ports.VehicleRepository
}

func NewCreateServiceOrder(
	repo ports.ServiceOrderRepository,
	customerRepo ports.CustomerRepository,
	vehicleRepo ports.VehicleRepository,
) *CreateServiceOrder {
	return &CreateServiceOrder{
		repo:         repo,
		customerRepo: customerRepo,
		vehicleRepo:  vehicleRepo,
	}
}

func (uc *CreateServiceOrder) Execute(ctx context.Context, input ports.CreateServiceOrderInput) (*entities.ServiceOrder, error) {
	customer, err := uc.customerRepo.FindByDocument(ctx, input.CustomerDocument)
	if err != nil {
		return nil, err
	}

	vehicle, err := uc.vehicleRepo.FindByPlate(ctx, input.VehiclePlate)
	if err != nil {
		return nil, err
	}

	if vehicle.CustomerID() != customer.ID() {
		return nil, domainerrors.ErrVehicleNotFromCustomer
	}

	so := entities.NewServiceOrder(customer.ID(), vehicle.ID())
	so.SetNotes(input.Notes)

	if err := uc.repo.Create(ctx, so); err != nil {
		return nil, err
	}

	return so, nil
}
