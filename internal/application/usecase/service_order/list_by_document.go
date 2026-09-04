package serviceorderuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports"
)

type ListServiceOrdersByDocument struct {
	repo         ports.ServiceOrderRepository
	requesterRepo ports.RequesterRepository
}

func NewListServiceOrdersByDocument(repo ports.ServiceOrderRepository, requesterRepo ports.RequesterRepository) *ListServiceOrdersByDocument {
	return &ListServiceOrdersByDocument{repo: repo, requesterRepo: requesterRepo}
}

func (uc *ListServiceOrdersByDocument) Execute(ctx context.Context, document string) ([]*entities.ServiceOrder, error) {
	requester, err := uc.requesterRepo.FindByDocument(ctx, document)
	if err != nil {
		return nil, err
	}

	return uc.repo.FindByRequesterID(ctx, requester.ID())
}
