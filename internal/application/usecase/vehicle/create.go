package vehicleuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
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
	// Placa ja validada na factory NewVehicle (fast-fail no dominio)

	// Verificar se o cliente existe
	_, err := uc.customerRepo.FindByID(ctx, v.CustomerID())
	if err != nil {
		return domainerrors.ErrNotFound
	}

	return uc.repo.Create(ctx, v)
}
