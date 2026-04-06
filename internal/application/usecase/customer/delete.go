package customeruc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/ports"
)

type DeleteCustomer struct {
	repo ports.CustomerRepository
}

func NewDeleteCustomer(repo ports.CustomerRepository) *DeleteCustomer {
	return &DeleteCustomer{repo: repo}
}

func (uc *DeleteCustomer) Execute(ctx context.Context, id string) error {
	// TODO: verificar dependencias (ordens de servico abertas)
	return uc.repo.Delete(ctx, id)
}
