package serviceorderuc

import (
	"context"
	"errors"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/ports"
)

type GetServiceOrderByCode struct {
	repo         ports.ServiceOrderRepository
	customerRepo ports.CustomerRepository
}

func NewGetServiceOrderByCode(repo ports.ServiceOrderRepository, customerRepo ports.CustomerRepository) *GetServiceOrderByCode {
	return &GetServiceOrderByCode{repo: repo, customerRepo: customerRepo}
}

func (uc *GetServiceOrderByCode) Execute(ctx context.Context, code int, customerCPF string) (*entities.ServiceOrder, error) {
	customer, err := uc.customerRepo.FindByDocument(ctx, customerCPF)
	if err != nil {
		if errors.Is(err, domainerrors.ErrNotFound) {
			return nil, domainerrors.ErrNotFound
		}
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
