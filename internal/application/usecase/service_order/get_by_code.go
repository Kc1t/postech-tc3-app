package serviceorderuc

import (
	"context"

	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports"
)

type GetServiceOrderByCode struct {
	repo         ports.ServiceOrderRepository
	requesterRepo ports.RequesterRepository
}

func NewGetServiceOrderByCode(repo ports.ServiceOrderRepository, requesterRepo ports.RequesterRepository) *GetServiceOrderByCode {
	return &GetServiceOrderByCode{repo: repo, requesterRepo: requesterRepo}
}

func (uc *GetServiceOrderByCode) Execute(ctx context.Context, code int, requesterDocument string) (*entities.ServiceOrder, error) {
	requester, err := uc.requesterRepo.FindByDocument(ctx, requesterDocument)
	if err != nil {
		return nil, err
	}

	so, err := uc.repo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if so.RequesterID() != requester.ID() {
		return nil, domainerrors.ErrNotFound
	}

	return so, nil
}
