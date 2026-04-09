package serviceuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/service"
	"github.com/fiap/postech-tc1/internal/ports"
)

type ListServices struct {
	repo ports.ServiceRepository
}

func NewListServices(repo ports.ServiceRepository) *ListServices {
	return &ListServices{repo: repo}
}

func (uc *ListServices) Execute(ctx context.Context) ([]*service.Service, error) {
	return uc.repo.FindAll(ctx)
}
