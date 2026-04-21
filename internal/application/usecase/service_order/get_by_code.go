package serviceorderuc

import (
	"context"

	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports"
)

type GetServiceOrderByCode struct {
	repo         ports.ServiceOrderRepository
	customerRepo ports.CustomerRepository
}

func NewGetServiceOrderByCode(repo ports.ServiceOrderRepository, customerRepo ports.CustomerRepository) *GetServiceOrderByCode {
	return &GetServiceOrderByCode{repo: repo, customerRepo: customerRepo}
}

func (uc *GetServiceOrderByCode) Execute(ctx context.Context, code int, customerDocument string) (*entities.ServiceOrder, error) {
	customer, err := uc.customerRepo.FindByDocument(ctx, customerDocument)
	if err != nil {
		return nil, err
	}

	so, err := uc.repo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if so.CustomerID() != customer.ID() {
		return nil, domainerrors.ErrNotFound
	}

	return so, nil
}
