package serviceorderuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports"
)

type UpdateServiceOrderStatus struct {
	repo ports.ServiceOrderRepository
}

func NewUpdateServiceOrderStatus(repo ports.ServiceOrderRepository) *UpdateServiceOrderStatus {
	return &UpdateServiceOrderStatus{repo: repo}
}

func (uc *UpdateServiceOrderStatus) Execute(ctx context.Context, id string, status entities.OrderStatus) error {
	// Buscar OS atual para validar transicao contra o status vigente
	so, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Maquina de estados valida a transicao no dominio (fast-fail)
	if err := so.UpdateStatus(status); err != nil {
		return err
	}

	return uc.repo.UpdateStatus(ctx, id, status)
}
