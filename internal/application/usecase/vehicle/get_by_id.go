package vehicleuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/vehicle"
	"github.com/fiap/postech-tc1/internal/ports"
)

type GetVehicle struct {
	repo ports.VehicleRepository
}

func NewGetVehicle(repo ports.VehicleRepository) *GetVehicle {
	return &GetVehicle{repo: repo}
}

func (uc *GetVehicle) Execute(ctx context.Context, id string) (*vehicle.Vehicle, error) {
	return uc.repo.FindByID(ctx, id)
}
