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

func (uc *CreateVehicle) Execute(ctx context.Context, customerDocument string, v *entities.Vehicle) error {
	customer, err := uc.customerRepo.FindByDocument(ctx, customerDocument)
	if err != nil {
		return err
	}

	v.SetCustomerID(customer.ID())
	return uc.repo.Create(ctx, v)
}
