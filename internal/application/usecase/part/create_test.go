package partuc

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestCreatePart_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockPartRepository(ctrl)
	repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	uc := NewCreatePart(repo)
	err := uc.Execute(context.Background(), entities.NewPart("Filtro", "Desc", "unidade", 49.90, 10))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestCreatePart_Execute_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoErr := errors.New("db error")
	repo := mocks.NewMockPartRepository(ctrl)
	repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(repoErr)

	uc := NewCreatePart(repo)
	if err := uc.Execute(context.Background(), entities.NewPart("Filtro", "Desc", "unidade", 49.90, 10)); !errors.Is(err, repoErr) {
		t.Fatalf("expected %v, got %v", repoErr, err)
	}
}
