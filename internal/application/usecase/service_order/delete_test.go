package serviceorderuc

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestDeleteServiceOrder_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	repo.EXPECT().Delete(gomock.Any(), "so-1").Return(nil)

	uc := NewDeleteServiceOrder(repo)
	if err := uc.Execute(context.Background(), "so-1"); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestDeleteServiceOrder_Execute_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoErr := errors.New("db error")
	repo := mocks.NewMockServiceOrderRepository(ctrl)
	repo.EXPECT().Delete(gomock.Any(), "so-1").Return(repoErr)

	uc := NewDeleteServiceOrder(repo)
	if err := uc.Execute(context.Background(), "so-1"); !errors.Is(err, repoErr) {
		t.Fatalf("expected %v, got %v", repoErr, err)
	}
}
