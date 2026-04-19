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
	part, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := part.AdjustStock(delta); err != nil {
		return err
	}

	return uc.repo.UpdateStock(ctx, id, delta)
}
