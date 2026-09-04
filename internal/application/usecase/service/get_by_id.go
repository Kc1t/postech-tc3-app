package serviceuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports"
)

type GetService struct {
	repo ports.ServiceRepository
}

func NewGetService(repo ports.ServiceRepository) *GetService {
	return &GetService{repo: repo}
}

func (uc *GetService) Execute(ctx context.Context, id string) (*entities.Service, error) {
	return uc.repo.FindByID(ctx, id)
}
