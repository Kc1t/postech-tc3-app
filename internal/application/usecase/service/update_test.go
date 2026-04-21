package serviceuc

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestUpdateService_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := entities.NewService(1, "Troca de óleo", "Desc", 150.0, 60)
	repo := mocks.NewMockServiceRepository(ctrl)
	repo.EXPECT().Update(gomock.Any(), s).Return(nil)

	uc := NewUpdateService(repo)
	if err := uc.Execute(context.Background(), s); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestUpdateService_Execute_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoErr := errors.New("db error")
	s := entities.NewService(1, "Troca de óleo", "Desc", 150.0, 60)
	repo := mocks.NewMockServiceRepository(ctrl)
	repo.EXPECT().Update(gomock.Any(), s).Return(repoErr)

	uc := NewUpdateService(repo)
	if err := uc.Execute(context.Background(), s); !errors.Is(err, repoErr) {
		t.Fatalf("expected %v, got %v", repoErr, err)
	}
}
