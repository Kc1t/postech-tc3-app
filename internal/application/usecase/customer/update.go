package customeruc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/customer"
	"github.com/fiap/postech-tc1/internal/ports"
)

type UpdateCustomer struct {
	repo ports.CustomerRepository
}

func NewUpdateCustomer(repo ports.CustomerRepository) *UpdateCustomer {
	return &UpdateCustomer{repo: repo}
}

func (uc *UpdateCustomer) Execute(ctx context.Context, c *customer.Customer) error {
	return uc.repo.Update(ctx, c)
}
