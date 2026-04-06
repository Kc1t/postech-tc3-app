package vehicleuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/vehicle"
	"github.com/fiap/postech-tc1/internal/ports"
)

type UpdateVehicle struct {
	repo ports.VehicleRepository
}

func NewUpdateVehicle(repo ports.VehicleRepository) *UpdateVehicle {
	return &UpdateVehicle{repo: repo}
}

func (uc *UpdateVehicle) Execute(ctx context.Context, v *vehicle.Vehicle) error {
	return uc.repo.Update(ctx, v)
}
