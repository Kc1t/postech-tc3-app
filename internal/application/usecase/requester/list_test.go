package requesteruc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestListRequesters_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := []*entities.Requester{
		entities.ReconstituteRequester("", "João", "123", "j@j.com", "11999", time.Now(), time.Now()),
		entities.ReconstituteRequester("", "Maria", "456", "m@m.com", "11888", time.Now(), time.Now()),
	}
	repo := mocks.NewMockRequesterRepository(ctrl)
	repo.EXPECT().FindAll(gomock.Any()).Return(expected, nil)

	uc := NewListRequesters(repo)
	got, err := uc.Execute(context.Background())
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 requesters, got %d", len(got))
	}
}

func TestListRequesters_Execute_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRequesterRepository(ctrl)
	repo.EXPECT().FindAll(gomock.Any()).Return([]*entities.Requester{}, nil)

	uc := NewListRequesters(repo)
	got, err := uc.Execute(context.Background())
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected 0 requesters, got %d", len(got))
	}
}

func TestListRequesters_Execute_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoErr := errors.New("db error")
	repo := mocks.NewMockRequesterRepository(ctrl)
	repo.EXPECT().FindAll(gomock.Any()).Return(nil, repoErr)

	uc := NewListRequesters(repo)
	_, err := uc.Execute(context.Background())
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected %v, got %v", repoErr, err)
	}
}
