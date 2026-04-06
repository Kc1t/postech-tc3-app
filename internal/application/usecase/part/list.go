package partuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/part"
	"github.com/fiap/postech-tc1/internal/ports"
)

type ListParts struct {
	repo ports.PartRepository
}

func NewListParts(repo ports.PartRepository) *ListParts {
	return &ListParts{repo: repo}
}

func (uc *ListParts) Execute(ctx context.Context) ([]*part.Part, error) {
	return uc.repo.FindAll(ctx)
}
