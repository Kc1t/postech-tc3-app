package serviceuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/ports"
)

type DeleteService struct {
	repo ports.ServiceRepository
}

func NewDeleteService(repo ports.ServiceRepository) *DeleteService {
	return &DeleteService{repo: repo}
}

func (uc *DeleteService) Execute(ctx context.Context, id string) error {
	// TODO: verificar se servico esta em uso em alguma OS ativa
	return uc.repo.Delete(ctx, id)
}
