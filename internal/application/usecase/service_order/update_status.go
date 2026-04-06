package serviceorderuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/serviceorder"
	"github.com/fiap/postech-tc1/internal/ports"
)

type UpdateServiceOrderStatus struct {
	repo ports.ServiceOrderRepository
}

func NewUpdateServiceOrderStatus(repo ports.ServiceOrderRepository) *UpdateServiceOrderStatus {
	return &UpdateServiceOrderStatus{repo: repo}
}

func (uc *UpdateServiceOrderStatus) Execute(ctx context.Context, id string, status serviceorder.Status) error {
	// TODO: validar transicao de status (maquina de estados)
	return uc.repo.UpdateStatus(ctx, id, status)
}
