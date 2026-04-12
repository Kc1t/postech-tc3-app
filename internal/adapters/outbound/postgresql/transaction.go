package postgresql

import (
	"context"

	"gorm.io/gorm"
)

type txKey struct{}

// GetDB retorna o *gorm.DB do contexto (transacao) ou o db padrao.
// Util quando um *gorm.DB transacional e injetado via contexto.
func GetDB(ctx context.Context, fallback *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return fallback
}
