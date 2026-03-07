package usecase

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/serviceorder"
	"github.com/fiap/postech-tc1/internal/ports"
)

type serviceOrderUseCase struct {
	repo         ports.ServiceOrderRepository
	customerRepo ports.CustomerRepository
	vehicleRepo  ports.VehicleRepository
	serviceRepo  ports.ServiceRepository
	partRepo     ports.PartRepository
}

func NewServiceOrderUseCase(
	repo ports.ServiceOrderRepository,
	customerRepo ports.CustomerRepository,
	vehicleRepo ports.VehicleRepository,
	serviceRepo ports.ServiceRepository,
	partRepo ports.PartRepository,
) ports.ServiceOrderUseCase {
	return &serviceOrderUseCase{
		repo:         repo,
		customerRepo: customerRepo,
		vehicleRepo:  vehicleRepo,
		serviceRepo:  serviceRepo,
		partRepo:     partRepo,
	}
}

func (uc *serviceOrderUseCase) Create(ctx context.Context, so *serviceorder.ServiceOrder) error {
	// TODO: validar cliente e veiculo, calcular total
	so.UpdateStatus(serviceorder.StatusReceived)
	return uc.repo.Create(ctx, so)
}

func (uc *serviceOrderUseCase) GetByID(ctx context.Context, id string) (*serviceorder.ServiceOrder, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *serviceOrderUseCase) GetAll(ctx context.Context) ([]*serviceorder.ServiceOrder, error) {
	return uc.repo.FindAll(ctx)
}

func (uc *serviceOrderUseCase) GetByCustomerID(ctx context.Context, customerID string) ([]*serviceorder.ServiceOrder, error) {
	return uc.repo.FindByCustomerID(ctx, customerID)
}

func (uc *serviceOrderUseCase) UpdateStatus(ctx context.Context, id string, status serviceorder.Status) error {
	// TODO: validar transicao de status (maquina de estados)
	return uc.repo.UpdateStatus(ctx, id, status)
}

func (uc *serviceOrderUseCase) Update(ctx context.Context, so *serviceorder.ServiceOrder) error {
	return uc.repo.Update(ctx, so)
}

func (uc *serviceOrderUseCase) Delete(ctx context.Context, id string) error {
	// TODO: apenas OS com status "received" podem ser deletadas
	return uc.repo.Delete(ctx, id)
}
