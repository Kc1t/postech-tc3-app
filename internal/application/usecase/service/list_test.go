package serviceuc

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestListServices_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := []*entities.Service{
		entities.NewService(1, "Troca de óleo", "Desc", 150.0, 60),
	}
	repo := mocks.NewMockServiceRepository(ctrl)
	repo.EXPECT().FindAll(gomock.Any()).Return(expected, nil)

	uc := NewListServices(repo)
	got, err := uc.Execute(context.Background())
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 service, got %d", len(got))
	}
}

func TestListServices_Execute_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoErr := errors.New("db error")
	repo := mocks.NewMockServiceRepository(ctrl)
	repo.EXPECT().FindAll(gomock.Any()).Return(nil, repoErr)

	uc := NewListServices(repo)
	_, err := uc.Execute(context.Background())
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected %v, got %v", repoErr, err)
	}
}
