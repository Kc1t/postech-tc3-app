package serviceorderuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/serviceorder"
	"github.com/fiap/postech-tc1/internal/ports"
)

type ListServiceOrders struct {
	repo ports.ServiceOrderRepository
}

func NewListServiceOrders(repo ports.ServiceOrderRepository) *ListServiceOrders {
	return &ListServiceOrders{repo: repo}
}

func (uc *ListServiceOrders) Execute(ctx context.Context) ([]*serviceorder.ServiceOrder, error) {
	return uc.repo.FindAll(ctx)
}
