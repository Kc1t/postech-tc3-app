package serviceuc

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestGetService_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := entities.NewService("Troca de óleo", "Desc", 150.0, 60)
	repo := mocks.NewMockServiceRepository(ctrl)
	repo.EXPECT().FindByID(gomock.Any(), "svc-1").Return(expected, nil)

	uc := NewGetService(repo)
	got, err := uc.Execute(context.Background(), "svc-1")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if got != expected {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

func TestGetService_Execute_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockServiceRepository(ctrl)
	repo.EXPECT().FindByID(gomock.Any(), "svc-x").Return(nil, domainerrors.ErrNotFound)

	uc := NewGetService(repo)
	_, err := uc.Execute(context.Background(), "svc-x")
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
