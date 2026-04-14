package customeruc

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestGetCustomer_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := entities.NewCustomer("João", "123", "j@j.com", "11999")
	repo := mocks.NewMockCustomerRepository(ctrl)
	repo.EXPECT().FindByID(gomock.Any(), "cust-1").Return(expected, nil)

	uc := NewGetCustomer(repo)
	got, err := uc.Execute(context.Background(), "cust-1")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if got != expected {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

func TestGetCustomer_Execute_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockCustomerRepository(ctrl)
	repo.EXPECT().FindByID(gomock.Any(), "cust-x").Return(nil, domainerrors.ErrNotFound)

	uc := NewGetCustomer(repo)
	_, err := uc.Execute(context.Background(), "cust-x")
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
