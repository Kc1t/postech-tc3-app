package serviceorderuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/ports"
	"github.com/fiap/postech-tc1/pkg/logger"
)

type UpdateServiceOrderStatusByCode struct {
	repo          ports.ServiceOrderRepository
	requesterRepo ports.RequesterRepository
	notifier      ports.EmailNotifier
}

func NewUpdateServiceOrderStatusByCode(
	repo ports.ServiceOrderRepository,
	requesterRepo ports.RequesterRepository,
	notifier ports.EmailNotifier,
) *UpdateServiceOrderStatusByCode {
	return &UpdateServiceOrderStatusByCode{repo: repo, requesterRepo: requesterRepo, notifier: notifier}
}

func (uc *UpdateServiceOrderStatusByCode) Execute(ctx context.Context, code int, requesterDocument string, newStatus entities.OrderStatus) error {
	requester, err := uc.requesterRepo.FindByDocument(ctx, requesterDocument)
	if err != nil {
		return err
	}

	so, err := uc.repo.FindByCode(ctx, code)
	if err != nil {
		return err
	}

	if so.RequesterID() != requester.ID() {
		return domainerrors.ErrNotFound
	}

	if err := so.AuthorizeRequesterTransition(newStatus); err != nil {
		return err
	}

	if newStatus == entities.StatusInExecution {
		if err := uc.repo.ApplyApprovalTransition(ctx, so); err != nil {
			return err
		}
	} else {
		if err := uc.repo.UpdateStatus(ctx, so.ID(), newStatus); err != nil {
			return err
		}
	}

	logger.FromContext(ctx).Info("service_order_status_changed",
		"service_order_code", so.Code(),
		"status", string(so.Status()),
		"actor", "requester")

	if uc.notifier != nil {
		if err := uc.notifier.NotifyStatusChange(ctx, requester.Email(), requester.Name(), so.Code(), so.Status()); err != nil {
			logger.FromContext(ctx).Error("integration_failure",
				"integration", "smtp", "stage", "send",
				"service_order_code", so.Code(), "error", err)
		}
	}

	return nil
}
