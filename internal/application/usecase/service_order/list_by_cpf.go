package serviceorderuc

import (
	"context"
	"errors"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/ports"
)

type ListServiceOrdersByCPF struct {
	repo         ports.ServiceOrderRepository
	customerRepo ports.CustomerRepository
}

func NewListServiceOrdersByCPF(repo ports.ServiceOrderRepository, customerRepo ports.CustomerRepository) *ListServiceOrdersByCPF {
	return &ListServiceOrdersByCPF{repo: repo, customerRepo: customerRepo}
}

func (uc *ListServiceOrdersByCPF) Execute(ctx context.Context, cpf string) ([]*entities.ServiceOrder, error) {
	customer, err := uc.customerRepo.FindByDocument(ctx, cpf)
	if err != nil {
		if errors.Is(err, domainerrors.ErrNotFound) {
			return nil, domainerrors.ErrNotFound
		}
		return nil, err
	}

	return uc.repo.FindByCustomerID(ctx, customer.ID())
}
