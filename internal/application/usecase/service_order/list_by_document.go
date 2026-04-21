package serviceorderuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports"
)

type ListServiceOrdersByDocument struct {
	repo         ports.ServiceOrderRepository
	customerRepo ports.CustomerRepository
}

func NewListServiceOrdersByDocument(repo ports.ServiceOrderRepository, customerRepo ports.CustomerRepository) *ListServiceOrdersByDocument {
	return &ListServiceOrdersByDocument{repo: repo, customerRepo: customerRepo}
}

func (uc *ListServiceOrdersByDocument) Execute(ctx context.Context, document string) ([]*entities.ServiceOrder, error) {
	customer, err := uc.customerRepo.FindByDocument(ctx, document)
	if err != nil {
		return nil, err
	}

	return uc.repo.FindByCustomerID(ctx, customer.ID())
}
