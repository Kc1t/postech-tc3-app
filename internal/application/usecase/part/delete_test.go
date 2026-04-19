package partuc

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestDeletePart_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockPartRepository(ctrl)
	repo.EXPECT().Delete(gomock.Any(), "part-1").Return(nil)

	uc := NewDeletePart(repo)
	if err := uc.Execute(context.Background(), "part-1"); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestDeletePart_Execute_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoErr := errors.New("db error")
	repo := mocks.NewMockPartRepository(ctrl)
	repo.EXPECT().Delete(gomock.Any(), "part-1").Return(repoErr)

	uc := NewDeletePart(repo)
	if err := uc.Execute(context.Background(), "part-1"); !errors.Is(err, repoErr) {
		t.Fatalf("expected %v, got %v", repoErr, err)
	}
}
