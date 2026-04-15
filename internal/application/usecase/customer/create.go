package customeruc

import (
	"context"
	"errors"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/ports"
)

type CreateCustomer struct {
	repo ports.CustomerRepository
}

func NewCreateCustomer(repo ports.CustomerRepository) *CreateCustomer {
	return &CreateCustomer{repo: repo}
}

func (uc *CreateCustomer) Execute(ctx context.Context, c *entities.Customer) error {
	// Documento ja validado na factory NewCustomer (fast-fail no dominio)

	// Verificar duplicidade de documento
	existing, err := uc.repo.FindByDocument(ctx, c.Document())
	if err != nil && !errors.Is(err, domainerrors.ErrNotFound) {
		return err
	}
	if existing != nil {
		return domainerrors.ErrAlreadyExists
	}

	return uc.repo.Create(ctx, c)
}
