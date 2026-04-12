package serviceuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports"
)

type ListServices struct {
	repo ports.ServiceRepository
}

func NewListServices(repo ports.ServiceRepository) *ListServices {
	return &ListServices{repo: repo}
}

func (uc *ListServices) Execute(ctx context.Context) ([]*entities.Service, error) {
	return uc.repo.FindAll(ctx)
}
