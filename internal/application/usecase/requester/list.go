package requesteruc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports"
)

type ListRequesters struct {
	repo ports.RequesterRepository
}

func NewListRequesters(repo ports.RequesterRepository) *ListRequesters {
	return &ListRequesters{repo: repo}
}

func (uc *ListRequesters) Execute(ctx context.Context) ([]*entities.Requester, error) {
	return uc.repo.FindAll(ctx)
}
