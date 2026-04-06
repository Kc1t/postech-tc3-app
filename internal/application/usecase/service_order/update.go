package serviceorderuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/serviceorder"
	"github.com/fiap/postech-tc1/internal/ports"
)

type UpdateServiceOrder struct {
	repo ports.ServiceOrderRepository
}

func NewUpdateServiceOrder(repo ports.ServiceOrderRepository) *UpdateServiceOrder {
	return &UpdateServiceOrder{repo: repo}
}

func (uc *UpdateServiceOrder) Execute(ctx context.Context, so *serviceorder.ServiceOrder) error {
	return uc.repo.Update(ctx, so)
}
