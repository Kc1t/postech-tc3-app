package serviceorderuc

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestUpdateServiceOrderStatus_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	repo.EXPECT().UpdateStatus(gomock.Any(), "so-1", entities.StatusInDiagnosis).Return(nil)

	uc := NewUpdateServiceOrderStatus(repo)
	if err := uc.Execute(context.Background(), "so-1", entities.StatusInDiagnosis); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestUpdateServiceOrderStatus_Execute_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoErr := errors.New("db error")
	repo := mocks.NewMockServiceOrderRepository(ctrl)
	repo.EXPECT().UpdateStatus(gomock.Any(), "so-1", entities.StatusFinished).Return(repoErr)

	uc := NewUpdateServiceOrderStatus(repo)
	if err := uc.Execute(context.Background(), "so-1", entities.StatusFinished); !errors.Is(err, repoErr) {
		t.Fatalf("expected %v, got %v", repoErr, err)
	}
}
