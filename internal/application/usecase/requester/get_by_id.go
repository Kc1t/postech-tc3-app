package requesteruc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports"
)

type GetRequester struct {
	repo ports.RequesterRepository
}

func NewGetRequester(repo ports.RequesterRepository) *GetRequester {
	return &GetRequester{repo: repo}
}

func (uc *GetRequester) Execute(ctx context.Context, id string) (*entities.Requester, error) {
	return uc.repo.FindByID(ctx, id)
}
