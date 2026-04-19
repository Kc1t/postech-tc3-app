package serviceorderuc

import (
	"context"
	"errors"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
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

func (uc *CreateServiceOrder) Execute(ctx context.Context, so *entities.ServiceOrder) error {
	// Validar existencia do cliente
	_, err := uc.customerRepo.FindByID(ctx, so.CustomerID())
	if err != nil {
		if errors.Is(err, domainerrors.ErrNotFound) {
			return domainerrors.ErrNotFound
		}
		return err
	}

	// Validar existencia do veiculo e vinculo com o cliente
	vehicle, err := uc.vehicleRepo.FindByID(ctx, so.VehicleID())
	if err != nil {
		if errors.Is(err, domainerrors.ErrNotFound) {
			return domainerrors.ErrNotFound
		}
		return err
	}
	if vehicle.CustomerID() != so.CustomerID() {
		return domainerrors.ErrVehicleNotFromCustomer
	}

	// Montar orcamento: buscar servicos em batch (1 query)
	if err := uc.buildServices(ctx, so); err != nil {
		return err
	}

	// Montar orcamento: buscar pecas em batch (1 query) e validar estoque
	if err := uc.buildParts(ctx, so); err != nil {
		return err
	}

	return uc.repo.Create(ctx, so)
}

func (uc *CreateServiceOrder) buildServices(ctx context.Context, so *entities.ServiceOrder) error {
	items := so.Services()
	if len(items) == 0 {
		return nil
	}

	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.ServiceID
	}

	services, err := uc.serviceRepo.FindByIDs(ctx, ids)
	if err != nil {
		return err
	}

	// Indexar por ID para lookup O(1)
	svcMap := make(map[string]*entities.Service, len(services))
	for _, svc := range services {
		svcMap[svc.ID()] = svc
	}

	// Validar que todos os servicos informados existem
	validated := make([]entities.ServiceItem, 0, len(items))
	for _, item := range items {
		svc, exists := svcMap[item.ServiceID]
		if !exists {
			return domainerrors.ErrNotFound
		}
		validated = append(validated, entities.ServiceItem{
			ServiceID:   svc.ID(),
			Description: svc.Name(),
			Price:       svc.Price(),
		})
	}

	so.SetServices(validated)
	return nil
}

func (uc *CreateServiceOrder) buildParts(ctx context.Context, so *entities.ServiceOrder) error {
	items := so.Parts()
	if len(items) == 0 {
		return nil
	}

	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.PartID
	}

	parts, err := uc.partRepo.FindByIDs(ctx, ids)
	if err != nil {
		return err
	}

	// Indexar por ID para lookup O(1)
	partMap := make(map[string]*entities.Part, len(parts))
	for _, p := range parts {
		partMap[p.ID()] = p
	}

	// Validar existencia e estoque de cada peca
	validated := make([]entities.PartItem, 0, len(items))
	for _, item := range items {
		part, exists := partMap[item.PartID]
		if !exists {
			return domainerrors.ErrNotFound
		}
		if part.Stock() < item.Quantity {
			return domainerrors.ErrInsufficientStock
		}
		validated = append(validated, entities.PartItem{
			PartID:      part.ID(),
			Description: part.Name(),
			Quantity:    item.Quantity,
			UnitPrice:   part.Price(),
		})
	}

	so.SetParts(validated)
	return nil
}
