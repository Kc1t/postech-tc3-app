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
	partRepo     ports.PartRepository
}

func NewUpdateServiceOrderStatusByCode(
	repo ports.ServiceOrderRepository,
	customerRepo ports.CustomerRepository,
	partRepo ports.PartRepository,
) *UpdateServiceOrderStatusByCode {
	return &UpdateServiceOrderStatusByCode{repo: repo, customerRepo: customerRepo, partRepo: partRepo}
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

	// Aprovacao pelo cliente (in_execution) baixa o estoque e persiste o
	// agregado completo — startedAt foi gravado pela entidade em UpdateStatus.
	// Recusa (received) nao toca timestamps, basta atualizar a coluna status.
	if newStatus == entities.StatusInExecution {
		if err := decrementStockForApproval(ctx, uc.partRepo, so.Parts()); err != nil {
			return err
		}
		return uc.repo.Update(ctx, so)
	}

	return uc.repo.UpdateStatus(ctx, so.ID(), newStatus)
}
