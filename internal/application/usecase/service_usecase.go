package usecase

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/service"
	"github.com/fiap/postech-tc1/internal/ports"
)

type serviceUseCase struct {
	repo ports.ServiceRepository
}

func NewServiceUseCase(repo ports.ServiceRepository) ports.ServiceUseCase {
	return &serviceUseCase{repo: repo}
}

func (uc *serviceUseCase) Create(ctx context.Context, s *service.Service) error {
	return uc.repo.Create(ctx, s)
}

func (uc *serviceUseCase) GetByID(ctx context.Context, id string) (*service.Service, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *serviceUseCase) GetAll(ctx context.Context) ([]*service.Service, error) {
	return uc.repo.FindAll(ctx)
}

func (uc *serviceUseCase) Update(ctx context.Context, s *service.Service) error {
	return uc.repo.Update(ctx, s)
}

func (uc *serviceUseCase) Delete(ctx context.Context, id string) error {
	// TODO: verificar se servico esta em uso em alguma OS ativa
	return uc.repo.Delete(ctx, id)
}
