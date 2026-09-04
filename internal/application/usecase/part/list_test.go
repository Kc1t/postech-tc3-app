package partuc

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestListParts_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := []*entities.Part{
		entities.NewPart("FAB-001", "Filtro", "Desc", "unidade", 49.90, 10),
	}
	repo := mocks.NewMockPartRepository(ctrl)
	repo.EXPECT().FindAll(gomock.Any()).Return(expected, nil)

	uc := NewListParts(repo)
	got, err := uc.Execute(context.Background())
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 part, got %d", len(got))
	}
}

func TestListParts_Execute_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoErr := errors.New("db error")
	repo := mocks.NewMockPartRepository(ctrl)
	repo.EXPECT().FindAll(gomock.Any()).Return(nil, repoErr)

	uc := NewListParts(repo)
	_, err := uc.Execute(context.Background())
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected %v, got %v", repoErr, err)
	}
}
