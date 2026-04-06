package customeruc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports"
)

type ListCustomers struct {
	repo ports.CustomerRepository
}

func NewListCustomers(repo ports.CustomerRepository) *ListCustomers {
	return &ListCustomers{repo: repo}
}

func (uc *ListCustomers) Execute(ctx context.Context) ([]*entities.Customer, error) {
	// TODO: paginacao
	return uc.repo.FindAll(ctx)
}
