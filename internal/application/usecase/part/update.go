package partuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/part"
	"github.com/fiap/postech-tc1/internal/ports"
)

type UpdatePart struct {
	repo ports.PartRepository
}

func NewUpdatePart(repo ports.PartRepository) *UpdatePart {
	return &UpdatePart{repo: repo}
}

func (uc *UpdatePart) Execute(ctx context.Context, p *part.Part) error {
	return uc.repo.Update(ctx, p)
}
