package serviceorderuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports"
)

type ListServiceOrdersByCustomer struct {
	repo ports.ServiceOrderRepository
}

func NewListServiceOrdersByCustomer(repo ports.ServiceOrderRepository) *ListServiceOrdersByCustomer {
	return &ListServiceOrdersByCustomer{repo: repo}
}

func (uc *ListServiceOrdersByCustomer) Execute(ctx context.Context, customerID string) ([]*entities.ServiceOrder, error) {
	return uc.repo.FindByCustomerID(ctx, customerID)
}
