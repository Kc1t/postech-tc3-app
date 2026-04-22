package vehicleuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports"
)

type CreateVehicle struct {
	repo         ports.VehicleRepository
	customerRepo ports.CustomerRepository
}

func NewCreateVehicle(repo ports.VehicleRepository, customerRepo ports.CustomerRepository) *CreateVehicle {
	return &CreateVehicle{repo: repo, customerRepo: customerRepo}
}

func (uc *CreateVehicle) Execute(ctx context.Context, input entities.VehicleInput) (*entities.Vehicle, error) {
	customer, err := uc.customerRepo.FindByDocument(ctx, input.CustomerDocument)
	if err != nil {
		return nil, err
	}

	vehicle, err := entities.NewVehicle(customer.ID(), input.Plate, input.Brand, input.Model, input.Year)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, vehicle); err != nil {
		return nil, err
	}

	return vehicle, nil
}
