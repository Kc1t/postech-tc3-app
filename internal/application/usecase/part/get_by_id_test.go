package partuc

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestGetPart_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := entities.NewPart("FAB-001", "Filtro", "Desc", "unidade", 49.90, 10)
	repo := mocks.NewMockPartRepository(ctrl)
	repo.EXPECT().FindByID(gomock.Any(), "part-1").Return(expected, nil)

	uc := NewGetPart(repo)
	got, err := uc.Execute(context.Background(), "part-1")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if got != expected {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

func TestGetPart_Execute_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockPartRepository(ctrl)
	repo.EXPECT().FindByID(gomock.Any(), "part-x").Return(nil, domainerrors.ErrNotFound)

	uc := NewGetPart(repo)
	_, err := uc.Execute(context.Background(), "part-x")
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
