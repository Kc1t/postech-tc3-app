package postgresql

import (
	"errors"
	"testing"

	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

func TestMapError_RecordNotFound(t *testing.T) {
	err := mapError(gorm.ErrRecordNotFound)
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestMapError_UniqueConstraint(t *testing.T) {
	pgErr := &pgconn.PgError{Code: "23505"}
	err := mapError(pgErr)
	if !errors.Is(err, domainerrors.ErrAlreadyExists) {
		t.Fatalf("expected ErrAlreadyExists, got %v", err)
	}
}

func TestMapError_OtherPgError(t *testing.T) {
	pgErr := &pgconn.PgError{Code: "23502"} // not null violation
	err := mapError(pgErr)
	if errors.Is(err, domainerrors.ErrNotFound) || errors.Is(err, domainerrors.ErrAlreadyExists) {
		t.Fatalf("expected original error to pass through, got domain error")
	}
	if err != pgErr {
		t.Fatalf("expected original error, got %v", err)
	}
}

func TestMapError_OtherError(t *testing.T) {
	original := errors.New("some other error")
	err := mapError(original)
	if err != original {
		t.Fatalf("expected original error, got %v", err)
	}
}
