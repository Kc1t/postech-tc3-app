package serviceorderuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/ports"
)

type UpdateServiceOrderStatus struct {
	repo        ports.ServiceOrderRepository
	serviceRepo ports.ServiceRepository
	partRepo    ports.PartRepository
}

func NewUpdateServiceOrderStatus(
	repo ports.ServiceOrderRepository,
	serviceRepo ports.ServiceRepository,
	partRepo ports.PartRepository,
) *UpdateServiceOrderStatus {
	return &UpdateServiceOrderStatus{
		repo:        repo,
		serviceRepo: serviceRepo,
		partRepo:    partRepo,
	}
}

func (uc *UpdateServiceOrderStatus) Execute(ctx context.Context, input entities.StatusUpdate) error {
	so, err := uc.repo.FindByID(ctx, input.ID)
	if err != nil {
		return err
	}

	// Maquina de estados valida a transicao no dominio (fast-fail)
	if err := so.UpdateStatus(input.Status); err != nil {
		return err
	}

	switch input.Status {
	case entities.StatusAwaitingApproval:
		// Monta o orcamento com servicos e pecas — precisa persistir o agregado completo.
		if err := uc.buildServices(ctx, so, input.ServiceCodes); err != nil {
			return err
		}
		if err := uc.buildParts(ctx, so, input.Parts); err != nil {
			return err
		}
		return uc.repo.Update(ctx, so)

	case entities.StatusInExecution:
		// Aprovacao do orcamento: decremento de estoque e persistencia da OS
		// rodam em uma unica transacao DB (rollback automatico em qualquer falha).
		// startedAt ja foi gravado pela entidade em UpdateStatus.
		return uc.repo.ApplyApprovalTransition(ctx, so)

	case entities.StatusFinished:
		// finishedAt foi gravado pela entidade — persistir agregado completo.
		return uc.repo.Update(ctx, so)

	default:
		// Transicoes sem side-effect em timestamps ou itens — basta atualizar a coluna status.
		return uc.repo.UpdateStatus(ctx, input.ID, input.Status)
	}
}

func (uc *UpdateServiceOrderStatus) buildServices(ctx context.Context, so *entities.ServiceOrder, codes []int) error {
	if len(codes) == 0 {
		return nil
	}

	services, err := uc.serviceRepo.FindByCodes(ctx, codes)
	if err != nil {
		return err
	}

	svcMap := make(map[int]*entities.Service, len(services))
	for _, svc := range services {
		svcMap[svc.Code()] = svc
	}

	items := make([]entities.ServiceItem, 0, len(codes))
	for _, code := range codes {
		svc, exists := svcMap[code]
		if !exists {
			return domainerrors.ErrNotFound
		}
		items = append(items, entities.ServiceItem{
			ServiceID:   svc.ID(),
			Description: svc.Name(),
			Price:       svc.Price(),
		})
	}

	so.SetServices(items)
	return nil
}

func (uc *UpdateServiceOrderStatus) buildParts(ctx context.Context, so *entities.ServiceOrder, partInputs []entities.OrderPartItem) error {
	if len(partInputs) == 0 {
		return nil
	}

	codes := make([]string, len(partInputs))
	for i, p := range partInputs {
		codes[i] = p.ManufacturerCode
	}

	parts, err := uc.partRepo.FindByManufacturerCodes(ctx, codes)
	if err != nil {
		return err
	}

	partMap := make(map[string]*entities.Part, len(parts))
	for _, p := range parts {
		partMap[p.ManufacturerCode()] = p
	}

	items := make([]entities.PartItem, 0, len(partInputs))
	for _, input := range partInputs {
		part, exists := partMap[input.ManufacturerCode]
		if !exists {
			return domainerrors.ErrNotFound
		}
		if part.Stock() < input.Quantity {
			return domainerrors.ErrInsufficientStock
		}
		items = append(items, entities.PartItem{
			PartID:      part.ID(),
			Description: part.Name(),
			Quantity:    input.Quantity,
			UnitPrice:   part.Price(),
		})
	}

	so.SetParts(items)
	return nil
}
