package customeruc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/customer"
	"github.com/fiap/postech-tc1/internal/ports"
)

type GetCustomer struct {
	repo ports.CustomerRepository
}

func NewGetCustomer(repo ports.CustomerRepository) *GetCustomer {
	return &GetCustomer{repo: repo}
}

func (uc *GetCustomer) Execute(ctx context.Context, id string) (*customer.Customer, error) {
	return uc.repo.FindByID(ctx, id)
}
