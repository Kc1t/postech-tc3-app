package customeruc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports"
)

type GetCustomer struct {
	repo ports.CustomerRepository
}

func NewGetCustomer(repo ports.CustomerRepository) *GetCustomer {
	return &GetCustomer{repo: repo}
}

func (uc *GetCustomer) Execute(ctx context.Context, id string) (*entities.Customer, error) {
	return uc.repo.FindByID(ctx, id)
}
