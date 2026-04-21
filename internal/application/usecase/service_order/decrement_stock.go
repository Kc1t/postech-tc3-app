package serviceorderuc

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports"
)

// decrementStockForApproval baixa o estoque das pecas reservadas na OS ao
// aprovar o orcamento (transicao awaiting_approval -> in_execution). Aplica
// duas camadas de atomicidade:
//   - A: cada UpdateStock e atomico no banco (SQL "WHERE stock + ? >= 0").
//   - B: se uma peca falhar, as pecas ja decrementadas sao revertidas em
//     memoria operacional via UpdateStock com delta positivo, e o erro
//     original propaga para impedir a transicao de status.
//
// Erros do rollback sao silenciados de proposito: se a reversao falhar (ex:
// DB indisponivel), nao ha acao melhor no escopo do MVP — o operador
// precisara conciliar o estoque manualmente. Transacao cross-aggregate
// (camada C) ficou explicitamente fora do escopo.
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
