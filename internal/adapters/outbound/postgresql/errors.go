package postgresql

import (
	"errors"

	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// mapError traduz erros de infraestrutura (GORM/PostgreSQL) para erros de domínio.
// Deve ser chamado em todos os métodos dos repositories antes de retornar um erro.
func mapError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domainerrors.ErrNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return domainerrors.ErrAlreadyExists
	}
	return err
}
