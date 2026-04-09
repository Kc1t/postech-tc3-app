package customeruc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/customer"
	"github.com/fiap/postech-tc1/internal/ports"
)

type GetCustomerByDocument struct {
	repo ports.CustomerRepository
}

func NewGetCustomerByDocument(repo ports.CustomerRepository) *GetCustomerByDocument {
	return &GetCustomerByDocument{repo: repo}
}

func (uc *GetCustomerByDocument) Execute(ctx context.Context, document string) (*customer.Customer, error) {
	return uc.repo.FindByDocument(ctx, document)
}
