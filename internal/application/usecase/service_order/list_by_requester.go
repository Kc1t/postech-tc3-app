package serviceorderuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports"
)

type ListServiceOrdersByRequester struct {
	repo ports.ServiceOrderRepository
}

func NewListServiceOrdersByRequester(repo ports.ServiceOrderRepository) *ListServiceOrdersByRequester {
	return &ListServiceOrdersByRequester{repo: repo}
}

func (uc *ListServiceOrdersByRequester) Execute(ctx context.Context, requesterID string) ([]*entities.ServiceOrder, error) {
	return uc.repo.FindByRequesterID(ctx, requesterID)
}
