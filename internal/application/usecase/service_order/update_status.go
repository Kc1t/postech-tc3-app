package serviceorderuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/ports"
	"github.com/fiap/postech-tc1/pkg/logger"
)

type UpdateServiceOrderStatus struct {
	repo          ports.ServiceOrderRepository
	serviceRepo   ports.ServiceRepository
	partRepo      ports.PartRepository
	requesterRepo ports.RequesterRepository
	notifier      ports.EmailNotifier
}

func NewUpdateServiceOrderStatus(
	repo ports.ServiceOrderRepository,
	serviceRepo ports.ServiceRepository,
	partRepo ports.PartRepository,
	requesterRepo ports.RequesterRepository,
	notifier ports.EmailNotifier,
) *UpdateServiceOrderStatus {
	return &UpdateServiceOrderStatus{
		repo:          repo,
		serviceRepo:   serviceRepo,
		partRepo:      partRepo,
		requesterRepo: requesterRepo,
		notifier:      notifier,
	}
}

func (uc *UpdateServiceOrderStatus) Execute(ctx context.Context, input entities.StatusUpdate) error {
	so, err := uc.repo.FindByID(ctx, input.ID)
	if err != nil {
		return err
	}

	if err := so.UpdateStatus(input.Status); err != nil {
		return err
	}

	switch input.Status {
	case entities.StatusAwaitingApproval:
		if err := uc.buildServices(ctx, so, input.ServiceCodes); err != nil {
			return err
		}
		if err := uc.buildParts(ctx, so, input.Parts); err != nil {
			return err
		}
		if err := uc.repo.Update(ctx, so); err != nil {
			return err
		}

	case entities.StatusInExecution:
		if err := uc.repo.ApplyApprovalTransition(ctx, so); err != nil {
			return err
		}

	case entities.StatusFinished:
		if err := uc.repo.Update(ctx, so); err != nil {
			return err
		}

	default:
		if err := uc.repo.UpdateStatus(ctx, input.ID, input.Status); err != nil {
			return err
		}
	}

	uc.logStatusChange(ctx, so, "staff")
	uc.notify(ctx, so)
	return nil
}

func (uc *UpdateServiceOrderStatus) logStatusChange(ctx context.Context, so *entities.ServiceOrder, actor string) {
	attrs := []any{
		"service_order_code", so.Code(),
		"status", string(so.Status()),
		"actor", actor,
	}

	if so.StartedAt() != nil && so.FinishedAt() != nil {
		attrs = append(attrs, "execution_seconds", so.FinishedAt().Sub(*so.StartedAt()).Seconds())
	}

	logger.FromContext(ctx).Info("service_order_status_changed", attrs...)
}

func (uc *UpdateServiceOrderStatus) notify(ctx context.Context, so *entities.ServiceOrder) {
	if uc.notifier == nil || uc.requesterRepo == nil {
		return
	}
	requester, err := uc.requesterRepo.FindByID(ctx, so.RequesterID())
	if err != nil {
		logger.FromContext(ctx).Error("integration_failure",
			"integration", "smtp", "stage", "requester_lookup",
			"service_order_code", so.Code(), "error", err)
		return
	}
	if err := uc.notifier.NotifyStatusChange(ctx, requester.Email(), requester.Name(), so.Code(), so.Status()); err != nil {
		logger.FromContext(ctx).Error("integration_failure",
			"integration", "smtp", "stage", "send",
			"service_order_code", so.Code(), "error", err)
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
