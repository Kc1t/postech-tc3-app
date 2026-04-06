package serviceorderuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/serviceorder"
	"github.com/fiap/postech-tc1/internal/ports"
)

type CreateServiceOrder struct {
	repo         ports.ServiceOrderRepository
	customerRepo ports.CustomerRepository
	vehicleRepo  ports.VehicleRepository
	serviceRepo  ports.ServiceRepository
	partRepo     ports.PartRepository
}

func NewCreateServiceOrder(
	repo ports.ServiceOrderRepository,
	customerRepo ports.CustomerRepository,
	vehicleRepo ports.VehicleRepository,
	serviceRepo ports.ServiceRepository,
	partRepo ports.PartRepository,
) *CreateServiceOrder {
	return &CreateServiceOrder{
		repo:         repo,
		customerRepo: customerRepo,
		vehicleRepo:  vehicleRepo,
		serviceRepo:  serviceRepo,
		partRepo:     partRepo,
	}
}

func (uc *CreateServiceOrder) Execute(ctx context.Context, so *serviceorder.ServiceOrder) error {
	// TODO: validar cliente e veiculo, calcular total
	so.UpdateStatus(serviceorder.StatusReceived)
	return uc.repo.Create(ctx, so)
}
