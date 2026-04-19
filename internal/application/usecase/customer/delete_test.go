package customeruc

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestDeleteCustomer_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockCustomerRepository(ctrl)
	repo.EXPECT().Delete(gomock.Any(), "cust-1").Return(nil)

	uc := NewDeleteCustomer(repo)
	if err := uc.Execute(context.Background(), "cust-1"); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestDeleteCustomer_Execute_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoErr := errors.New("db error")
	repo := mocks.NewMockCustomerRepository(ctrl)
	repo.EXPECT().Delete(gomock.Any(), "cust-1").Return(repoErr)

	uc := NewDeleteCustomer(repo)
	if err := uc.Execute(context.Background(), "cust-1"); !errors.Is(err, repoErr) {
		t.Fatalf("expected %v, got %v", repoErr, err)
	}
}
