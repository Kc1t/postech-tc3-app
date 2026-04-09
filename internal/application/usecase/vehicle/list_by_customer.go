package vehicleuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/vehicle"
	"github.com/fiap/postech-tc1/internal/ports"
)

type ListVehiclesByCustomer struct {
	repo ports.VehicleRepository
}

func NewListVehiclesByCustomer(repo ports.VehicleRepository) *ListVehiclesByCustomer {
	return &ListVehiclesByCustomer{repo: repo}
}

func (uc *ListVehiclesByCustomer) Execute(ctx context.Context, customerID string) ([]*vehicle.Vehicle, error) {
	return uc.repo.FindByCustomerID(ctx, customerID)
}
