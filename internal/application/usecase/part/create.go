package partuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports"
)

type CreatePart struct {
	repo ports.PartRepository
}

func NewCreatePart(repo ports.PartRepository) *CreatePart {
	return &CreatePart{repo: repo}
}

func (uc *CreatePart) Execute(ctx context.Context, p *entities.Part) error {
	return uc.repo.Create(ctx, p)
}
