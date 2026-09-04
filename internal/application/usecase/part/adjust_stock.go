package partuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/ports"
)

type AdjustPartStock struct {
	repo ports.PartRepository
}

func NewAdjustPartStock(repo ports.PartRepository) *AdjustPartStock {
	return &AdjustPartStock{repo: repo}
}

func (uc *AdjustPartStock) Execute(ctx context.Context, id string, delta int) error {
	return uc.repo.UpdateStock(ctx, id, delta)
}
