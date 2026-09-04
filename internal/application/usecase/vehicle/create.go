package vehicleuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports"
)

type CreateVehicle struct {
	repo         ports.VehicleRepository
	requesterRepo ports.RequesterRepository
}

func NewCreateVehicle(repo ports.VehicleRepository, requesterRepo ports.RequesterRepository) *CreateVehicle {
	return &CreateVehicle{repo: repo, requesterRepo: requesterRepo}
}

func (uc *CreateVehicle) Execute(ctx context.Context, input entities.VehicleInput) (*entities.Vehicle, error) {
	requester, err := uc.requesterRepo.FindByDocument(ctx, input.RequesterDocument)
	if err != nil {
		return nil, err
	}

	vehicle, err := entities.NewVehicle(requester.ID(), input.Plate, input.Brand, input.Model, input.Year)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, vehicle); err != nil {
		return nil, err
	}

	return vehicle, nil
}
