package serviceuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/service"
	"github.com/fiap/postech-tc1/internal/ports"
)

type CreateService struct {
	repo ports.ServiceRepository
}

func NewCreateService(repo ports.ServiceRepository) *CreateService {
	return &CreateService{repo: repo}
}

func (uc *CreateService) Execute(ctx context.Context, s *service.Service) error {
	return uc.repo.Create(ctx, s)
}
