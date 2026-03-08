package usecase

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/vehicle"
	"github.com/fiap/postech-tc1/internal/ports"
)

type vehicleUseCase struct {
	repo         ports.VehicleRepository
	customerRepo ports.CustomerRepository
}

func NewVehicleUseCase(repo ports.VehicleRepository, customerRepo ports.CustomerRepository) ports.VehicleUseCase {
	return &vehicleUseCase{repo: repo, customerRepo: customerRepo}
}

func (uc *vehicleUseCase) Create(ctx context.Context, v *vehicle.Vehicle) error {
	// TODO: validar placa, verificar se cliente existe
	return uc.repo.Create(ctx, v)
}

func (uc *vehicleUseCase) GetByID(ctx context.Context, id string) (*vehicle.Vehicle, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *vehicleUseCase) GetByCustomerID(ctx context.Context, customerID string) ([]*vehicle.Vehicle, error) {
	return uc.repo.FindByCustomerID(ctx, customerID)
}

func (uc *vehicleUseCase) GetAll(ctx context.Context) ([]*vehicle.Vehicle, error) {
	return uc.repo.FindAll(ctx)
}

func (uc *vehicleUseCase) Update(ctx context.Context, v *vehicle.Vehicle) error {
	return uc.repo.Update(ctx, v)
}

func (uc *vehicleUseCase) Delete(ctx context.Context, id string) error {
	// TODO: verificar OS abertas para o veiculo
	return uc.repo.Delete(ctx, id)
}
