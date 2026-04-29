package vehicleuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports"
)

type ListVehiclesByRequester struct {
	repo ports.VehicleRepository
}

func NewListVehiclesByRequester(repo ports.VehicleRepository) *ListVehiclesByRequester {
	return &ListVehiclesByRequester{repo: repo}
}

func (uc *ListVehiclesByRequester) Execute(ctx context.Context, requesterID string) ([]*entities.Vehicle, error) {
	return uc.repo.FindByRequesterID(ctx, requesterID)
}
