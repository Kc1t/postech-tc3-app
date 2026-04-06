package vehicleuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/vehicle"
	"github.com/fiap/postech-tc1/internal/ports"
)

type ListVehicles struct {
	repo ports.VehicleRepository
}

func NewListVehicles(repo ports.VehicleRepository) *ListVehicles {
	return &ListVehicles{repo: repo}
}

func (uc *ListVehicles) Execute(ctx context.Context) ([]*vehicle.Vehicle, error) {
	return uc.repo.FindAll(ctx)
}
