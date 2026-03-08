package usecase

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/customer"
	"github.com/fiap/postech-tc1/internal/ports"
)

type customerUseCase struct {
	repo ports.CustomerRepository
}

func NewCustomerUseCase(repo ports.CustomerRepository) ports.CustomerUseCase {
	return &customerUseCase{repo: repo}
}

func (uc *customerUseCase) Create(ctx context.Context, c *customer.Customer) error {
	// TODO: validar CPF/CNPJ, verificar duplicidade
	return uc.repo.Create(ctx, c)
}

func (uc *customerUseCase) GetByID(ctx context.Context, id string) (*customer.Customer, error) {
	// TODO: implementar
	return uc.repo.FindByID(ctx, id)
}

func (uc *customerUseCase) GetByDocument(ctx context.Context, document string) (*customer.Customer, error) {
	// TODO: implementar
	return uc.repo.FindByDocument(ctx, document)
}

func (uc *customerUseCase) GetAll(ctx context.Context) ([]*customer.Customer, error) {
	// TODO: paginacao
	return uc.repo.FindAll(ctx)
}

func (uc *customerUseCase) Update(ctx context.Context, c *customer.Customer) error {
	// TODO: implementar
	return uc.repo.Update(ctx, c)
}

func (uc *customerUseCase) Delete(ctx context.Context, id string) error {
	// TODO: verificar dependencias (ordens de servico abertas)
	return uc.repo.Delete(ctx, id)
}
