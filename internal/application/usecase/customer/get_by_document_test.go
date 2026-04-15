package customeruc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestGetCustomerByDocument_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := entities.ReconstituteCustomer("", "João", "12345678901", "j@j.com", "11999", time.Now(), time.Now())
	repo := mocks.NewMockCustomerRepository(ctrl)
	repo.EXPECT().FindByDocument(gomock.Any(), "12345678901").Return(expected, nil)

	uc := NewGetCustomerByDocument(repo)
	got, err := uc.Execute(context.Background(), "12345678901")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if got != expected {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

func TestGetCustomerByDocument_Execute_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockCustomerRepository(ctrl)
	repo.EXPECT().FindByDocument(gomock.Any(), "00000000000").Return(nil, domainerrors.ErrNotFound)

	uc := NewGetCustomerByDocument(repo)
	_, err := uc.Execute(context.Background(), "00000000000")
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
