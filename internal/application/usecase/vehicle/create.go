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

func (uc *CreateVehicle) Execute(ctx context.Context, v *entities.Vehicle) error {
	// TODO: validar placa, verificar se cliente existe
	return uc.repo.Create(ctx, v)
}
