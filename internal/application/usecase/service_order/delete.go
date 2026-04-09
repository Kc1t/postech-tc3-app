package serviceorderuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/ports"
)

type DeleteServiceOrder struct {
	repo ports.ServiceOrderRepository
}

func NewDeleteServiceOrder(repo ports.ServiceOrderRepository) *DeleteServiceOrder {
	return &DeleteServiceOrder{repo: repo}
}

func (uc *DeleteServiceOrder) Execute(ctx context.Context, id string) error {
	// TODO: apenas OS com status "received" podem ser deletadas
	return uc.repo.Delete(ctx, id)
}
