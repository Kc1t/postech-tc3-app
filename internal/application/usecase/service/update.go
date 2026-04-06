package serviceuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/service"
	"github.com/fiap/postech-tc1/internal/ports"
)

type UpdateService struct {
	repo ports.ServiceRepository
}

func NewUpdateService(repo ports.ServiceRepository) *UpdateService {
	return &UpdateService{repo: repo}
}

func (uc *UpdateService) Execute(ctx context.Context, s *service.Service) error {
	return uc.repo.Update(ctx, s)
}
