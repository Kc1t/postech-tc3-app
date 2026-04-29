package requesteruc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports"
)

type UpdateRequester struct {
	repo ports.RequesterRepository
}

func NewUpdateRequester(repo ports.RequesterRepository) *UpdateRequester {
	return &UpdateRequester{repo: repo}
}

func (uc *UpdateRequester) Execute(ctx context.Context, c *entities.Requester) error {
	return uc.repo.Update(ctx, c)
}
