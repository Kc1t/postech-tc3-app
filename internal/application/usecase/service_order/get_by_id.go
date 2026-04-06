package serviceorderuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/serviceorder"
	"github.com/fiap/postech-tc1/internal/ports"
)

type GetServiceOrder struct {
	repo ports.ServiceOrderRepository
}

func NewGetServiceOrder(repo ports.ServiceOrderRepository) *GetServiceOrder {
	return &GetServiceOrder{repo: repo}
}

func (uc *GetServiceOrder) Execute(ctx context.Context, id string) (*serviceorder.ServiceOrder, error) {
	return uc.repo.FindByID(ctx, id)
}
