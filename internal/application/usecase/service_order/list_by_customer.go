package serviceorderuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/serviceorder"
	"github.com/fiap/postech-tc1/internal/ports"
)

type ListServiceOrdersByCustomer struct {
	repo ports.ServiceOrderRepository
}

func NewListServiceOrdersByCustomer(repo ports.ServiceOrderRepository) *ListServiceOrdersByCustomer {
	return &ListServiceOrdersByCustomer{repo: repo}
}

func (uc *ListServiceOrdersByCustomer) Execute(ctx context.Context, customerID string) ([]*serviceorder.ServiceOrder, error) {
	return uc.repo.FindByCustomerID(ctx, customerID)
}
