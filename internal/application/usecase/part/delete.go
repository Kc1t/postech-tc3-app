package partuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/ports"
)

type DeletePart struct {
	repo ports.PartRepository
}

func NewDeletePart(repo ports.PartRepository) *DeletePart {
	return &DeletePart{repo: repo}
}

func (uc *DeletePart) Execute(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}
