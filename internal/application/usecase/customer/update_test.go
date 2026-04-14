package customeruc

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestUpdateCustomer_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := entities.NewCustomer("João", "123", "j@j.com", "11999")
	repo := mocks.NewMockCustomerRepository(ctrl)
	repo.EXPECT().Update(gomock.Any(), c).Return(nil)

	uc := NewUpdateCustomer(repo)
	if err := uc.Execute(context.Background(), c); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestUpdateCustomer_Execute_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoErr := errors.New("db error")
	c := entities.NewCustomer("João", "123", "j@j.com", "11999")
	repo := mocks.NewMockCustomerRepository(ctrl)
	repo.EXPECT().Update(gomock.Any(), c).Return(repoErr)

	uc := NewUpdateCustomer(repo)
	if err := uc.Execute(context.Background(), c); !errors.Is(err, repoErr) {
		t.Fatalf("expected %v, got %v", repoErr, err)
	}
}
