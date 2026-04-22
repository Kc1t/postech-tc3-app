package serviceorderuc

import (
	"context"

	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports"
)

type UpdateServiceOrderStatusByCode struct {
	repo         ports.ServiceOrderRepository
	customerRepo ports.CustomerRepository
}

func NewUpdateServiceOrderStatusByCode(
	repo ports.ServiceOrderRepository,
	customerRepo ports.CustomerRepository,
) *UpdateServiceOrderStatusByCode {
	return &UpdateServiceOrderStatusByCode{repo: repo, customerRepo: customerRepo}
}

func (uc *UpdateServiceOrderStatusByCode) Execute(ctx context.Context, code int, customerDocument string, newStatus entities.OrderStatus) error {
	customer, err := uc.customerRepo.FindByDocument(ctx, customerDocument)
	if err != nil {
		return err
	}

	so, err := uc.repo.FindByCode(ctx, code)
	if err != nil {
		return err
	}

	if so.CustomerID() != customer.ID() {
		return domainerrors.ErrNotFound
	}

	if err := so.AuthorizeCustomerTransition(newStatus); err != nil {
		return err
	}

	// Aprovacao pelo cliente (in_execution): decremento de estoque + persistencia
	// da OS numa unica transacao DB. Recusa (received) nao toca timestamps,
	// basta atualizar a coluna status.
	if newStatus == entities.StatusInExecution {
		return uc.repo.ApplyApprovalTransition(ctx, so)
	}

	return uc.repo.UpdateStatus(ctx, so.ID(), newStatus)
}
