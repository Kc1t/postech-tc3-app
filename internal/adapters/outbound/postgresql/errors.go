package postgresql

import (
	"errors"

	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

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
