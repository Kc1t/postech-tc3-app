package serviceorderuc

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestGetServiceOrder_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := entities.NewServiceOrder("cust-1", "veh-1")
	repo := mocks.NewMockServiceOrderRepository(ctrl)
	repo.EXPECT().FindByID(gomock.Any(), "so-1").Return(expected, nil)

	uc := NewGetServiceOrder(repo)
	got, err := uc.Execute(context.Background(), "so-1")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if got != expected {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

func TestGetServiceOrder_Execute_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	repo.EXPECT().FindByID(gomock.Any(), "so-x").Return(nil, domainerrors.ErrNotFound)

	uc := NewGetServiceOrder(repo)
	_, err := uc.Execute(context.Background(), "so-x")
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
