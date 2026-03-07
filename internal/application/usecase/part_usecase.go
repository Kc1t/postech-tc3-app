package usecase

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/part"
	"github.com/fiap/postech-tc1/internal/ports"
)

type partUseCase struct {
	repo ports.PartRepository
}

func NewPartUseCase(repo ports.PartRepository) ports.PartUseCase {
	return &partUseCase{repo: repo}
}

func (uc *partUseCase) Create(ctx context.Context, p *part.Part) error {
	return uc.repo.Create(ctx, p)
}

func (uc *partUseCase) GetByID(ctx context.Context, id string) (*part.Part, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *partUseCase) GetAll(ctx context.Context) ([]*part.Part, error) {
	return uc.repo.FindAll(ctx)
}

func (uc *partUseCase) Update(ctx context.Context, p *part.Part) error {
	return uc.repo.Update(ctx, p)
}

func (uc *partUseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *partUseCase) AdjustStock(ctx context.Context, id string, delta int) error {
	// TODO: nao permitir estoque negativo
	return uc.repo.UpdateStock(ctx, id, delta)
}
