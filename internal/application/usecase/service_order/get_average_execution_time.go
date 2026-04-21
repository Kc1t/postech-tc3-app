package serviceorderuc

import (
	"context"
	"time"

	"github.com/fiap/postech-tc1/internal/ports"
)

// GetAverageExecutionTime expoe a metrica de tempo medio de execucao dos
// servicos — calculada como a media do intervalo entre aprovacao do orcamento
// (startedAt) e finalizacao do servico (finishedAt).
type GetAverageExecutionTime struct {
	repo ports.ServiceOrderRepository
}

func NewGetAverageExecutionTime(repo ports.ServiceOrderRepository) *GetAverageExecutionTime {
	return &GetAverageExecutionTime{repo: repo}
}

func (uc *GetAverageExecutionTime) Execute(ctx context.Context) (time.Duration, error) {
	return uc.repo.AverageExecutionTime(ctx)
}
