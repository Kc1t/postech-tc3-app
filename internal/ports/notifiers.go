package ports

import (
	"context"

	"github.com/fiap/postech-tc1/internal/domain/entities"
)

//go:generate mockgen -source=./notifiers.go -destination=./mocks/notifiers.go -package=mocks

// EmailNotifier envia notificacoes por e-mail ao cliente quando o status da OS muda.
type EmailNotifier interface {
	NotifyStatusChange(ctx context.Context, toEmail, toName string, orderCode int, status entities.OrderStatus) error
}
