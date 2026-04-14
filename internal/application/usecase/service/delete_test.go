package serviceuc

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestDeleteService_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockServiceRepository(ctrl)
	repo.EXPECT().Delete(gomock.Any(), "svc-1").Return(nil)

	uc := NewDeleteService(repo)
	if err := uc.Execute(context.Background(), "svc-1"); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestDeleteService_Execute_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoErr := errors.New("db error")
	repo := mocks.NewMockServiceRepository(ctrl)
	repo.EXPECT().Delete(gomock.Any(), "svc-1").Return(repoErr)

	uc := NewDeleteService(repo)
	if err := uc.Execute(context.Background(), "svc-1"); !errors.Is(err, repoErr) {
		t.Fatalf("expected %v, got %v", repoErr, err)
	}
}
