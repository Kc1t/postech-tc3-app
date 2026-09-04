package partuc

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestUpdatePart_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	p := entities.NewPart("FAB-001", "Filtro", "Desc", "unidade", 49.90, 10)
	repo := mocks.NewMockPartRepository(ctrl)
	repo.EXPECT().Update(gomock.Any(), p).Return(nil)

	uc := NewUpdatePart(repo)
	if err := uc.Execute(context.Background(), p); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestUpdatePart_Execute_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoErr := errors.New("db error")
	p := entities.NewPart("FAB-001", "Filtro", "Desc", "unidade", 49.90, 10)
	repo := mocks.NewMockPartRepository(ctrl)
	repo.EXPECT().Update(gomock.Any(), p).Return(repoErr)

	uc := NewUpdatePart(repo)
	if err := uc.Execute(context.Background(), p); !errors.Is(err, repoErr) {
		t.Fatalf("expected %v, got %v", repoErr, err)
	}
}
