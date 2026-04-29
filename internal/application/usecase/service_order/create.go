package serviceorderuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/ports"
)

type CreateServiceOrder struct {
	repo         ports.ServiceOrderRepository
	requesterRepo ports.RequesterRepository
	vehicleRepo  ports.VehicleRepository
}

func NewCreateServiceOrder(
	repo ports.ServiceOrderRepository,
	requesterRepo ports.RequesterRepository,
	vehicleRepo ports.VehicleRepository,
) *CreateServiceOrder {
	return &CreateServiceOrder{
		repo:         repo,
		requesterRepo: requesterRepo,
		vehicleRepo:  vehicleRepo,
	}
}

func (uc *CreateServiceOrder) Execute(ctx context.Context, input entities.ServiceOrderInput) (*entities.ServiceOrder, error) {
	requester, err := uc.requesterRepo.FindByDocument(ctx, input.RequesterDocument)
	if err != nil {
		return nil, err
	}

	vehicle, err := uc.vehicleRepo.FindByPlate(ctx, input.VehiclePlate)
	if err != nil {
		return nil, err
	}

	if vehicle.RequesterID() != requester.ID() {
		return nil, domainerrors.ErrVehicleNotFromRequester
	}

	so := entities.NewServiceOrder(requester.ID(), vehicle.ID())
	so.SetNotes(input.Notes)

	if err := uc.repo.Create(ctx, so); err != nil {
		return nil, err
	}

	return so, nil
}
