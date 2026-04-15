package serviceorderuc

import (
	"context"
	"errors"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/ports"
)

type ApproveServiceOrder struct {
	repo         ports.ServiceOrderRepository
	customerRepo ports.CustomerRepository
}

func NewApproveServiceOrder(repo ports.ServiceOrderRepository, customerRepo ports.CustomerRepository) *ApproveServiceOrder {
	return &ApproveServiceOrder{repo: repo, customerRepo: customerRepo}
}

func (uc *ApproveServiceOrder) Execute(ctx context.Context, code int, customerCPF string) error {
	customer, err := uc.customerRepo.FindByDocument(ctx, customerCPF)
	if err != nil {
		if errors.Is(err, domainerrors.ErrNotFound) {
			return domainerrors.ErrNotFound
		}
		return err
	}

	so, err := uc.repo.FindByCode(ctx, code)
	if err != nil {
		return err
	}

	if so.CustomerID() != customer.ID() {
		return domainerrors.ErrNotFound
	}

	// awaiting_approval → in_execution
	if err := so.UpdateStatus(entities.StatusInExecution); err != nil {
		return err
	}

	return uc.repo.UpdateStatus(ctx, so.ID(), entities.StatusInExecution)
}
