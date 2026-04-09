package customeruc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/customer"
	"github.com/fiap/postech-tc1/internal/ports"
)

type CreateCustomer struct {
	repo ports.CustomerRepository
}

func NewCreateCustomer(repo ports.CustomerRepository) *CreateCustomer {
	return &CreateCustomer{repo: repo}
}

func (uc *CreateCustomer) Execute(ctx context.Context, c *customer.Customer) error {
	// TODO: validar CPF/CNPJ, verificar duplicidade
	return uc.repo.Create(ctx, c)
}
