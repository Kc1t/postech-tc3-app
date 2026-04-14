package customeruc

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestCreateCustomer_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockCustomerRepository(ctrl)
	repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	uc := NewCreateCustomer(repo)
	err := uc.Execute(context.Background(), entities.NewCustomer("João", "12345678901", "j@j.com", "11999"))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestCreateCustomer_Execute_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoErr := errors.New("db error")
	repo := mocks.NewMockCustomerRepository(ctrl)
	repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(repoErr)

	uc := NewCreateCustomer(repo)
	err := uc.Execute(context.Background(), entities.NewCustomer("João", "12345678901", "j@j.com", "11999"))
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected %v, got %v", repoErr, err)
	}
}
