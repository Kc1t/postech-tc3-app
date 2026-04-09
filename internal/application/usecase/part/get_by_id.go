package partuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/part"
	"github.com/fiap/postech-tc1/internal/ports"
)

type GetPart struct {
	repo ports.PartRepository
}

func NewGetPart(repo ports.PartRepository) *GetPart {
	return &GetPart{repo: repo}
}

func (uc *GetPart) Execute(ctx context.Context, id string) (*part.Part, error) {
	return uc.repo.FindByID(ctx, id)
}
