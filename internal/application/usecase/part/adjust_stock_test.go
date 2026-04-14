package partuc

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestAdjustPartStock_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockPartRepository(ctrl)
	repo.EXPECT().UpdateStock(gomock.Any(), "part-1", 5).Return(nil)

	uc := NewAdjustPartStock(repo)
	if err := uc.Execute(context.Background(), "part-1", 5); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestAdjustPartStock_Execute_Negative(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockPartRepository(ctrl)
	repo.EXPECT().UpdateStock(gomock.Any(), "part-1", -3).Return(nil)

	uc := NewAdjustPartStock(repo)
	if err := uc.Execute(context.Background(), "part-1", -3); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestAdjustPartStock_Execute_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoErr := errors.New("db error")
	repo := mocks.NewMockPartRepository(ctrl)
	repo.EXPECT().UpdateStock(gomock.Any(), "part-1", 5).Return(repoErr)

	uc := NewAdjustPartStock(repo)
	if err := uc.Execute(context.Background(), "part-1", 5); !errors.Is(err, repoErr) {
		t.Fatalf("expected %v, got %v", repoErr, err)
	}
}
