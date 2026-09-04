package requesteruc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports"
)

type GetRequesterByDocument struct {
	repo ports.RequesterRepository
}

func NewGetRequesterByDocument(repo ports.RequesterRepository) *GetRequesterByDocument {
	return &GetRequesterByDocument{repo: repo}
}

func (uc *GetRequesterByDocument) Execute(ctx context.Context, document string) (*entities.Requester, error) {
	return uc.repo.FindByDocument(ctx, document)
}
