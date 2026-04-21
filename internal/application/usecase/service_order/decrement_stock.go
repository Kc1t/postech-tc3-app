package serviceorderuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports"
)

// decrementStockForApproval baixa o estoque das pecas ao aprovar o
// orcamento (awaiting_approval -> in_execution).
//
// Atomicidade por peca: SQL "stock + ? >= 0". Entre pecas: rollback manual
// (erro do rollback e silenciado). Cross-aggregate (estoque + status da OS):
// fora de escopo no MVP.
//
// TODO N+1: ate 2N-1 round-trips em caso de falha parcial. Solucao prevista:
// metodo em lote no PartRepository com transacao local + UPDATE CASE WHEN.
func decrementStockForApproval(ctx context.Context, partRepo ports.PartRepository, parts []entities.PartItem) error {
	decremented := make([]entities.PartItem, 0, len(parts))
	for _, p := range parts {
		if err := partRepo.UpdateStock(ctx, p.PartID, -p.Quantity); err != nil {
			for _, r := range decremented {
				_ = partRepo.UpdateStock(ctx, r.PartID, r.Quantity)
			}
			return err
		}
		decremented = append(decremented, p)
	}
	return nil
}
