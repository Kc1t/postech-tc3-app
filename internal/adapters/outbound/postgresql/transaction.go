package postgresql

import (
	"context"

	"github.com/fiap/postech-tc1/internal/ports"
	"gorm.io/gorm"
)

type transactionManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) ports.TransactionManager {
	return &transactionManager{db: db}
}

func (tm *transactionManager) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return tm.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Passa o tx via context pra que os repos possam usar
		txCtx := context.WithValue(ctx, txKey{}, tx)
		return fn(txCtx)
	})
}

type txKey struct{}

// GetDB retorna o *gorm.DB do contexto (transacao) ou o db padrao.
func GetDB(ctx context.Context, fallback *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return fallback
}
