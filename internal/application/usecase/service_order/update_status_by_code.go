package serviceorderuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/ports"
)

type UpdateServiceOrderStatusByCode struct {
	repo          ports.ServiceOrderRepository
	requesterRepo ports.RequesterRepository
}

func NewUpdateServiceOrderStatusByCode(
	repo ports.ServiceOrderRepository,
	requesterRepo ports.RequesterRepository,
) *UpdateServiceOrderStatusByCode {
	return &UpdateServiceOrderStatusByCode{repo: repo, requesterRepo: requesterRepo}
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
		return uc.repo.ApplyApprovalTransition(ctx, so)
	}

	return uc.repo.UpdateStatus(ctx, so.ID(), newStatus)
}
