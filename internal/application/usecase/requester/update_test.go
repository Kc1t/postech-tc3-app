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

func TestUpdateRequester_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := entities.ReconstituteRequester("", "João", "123", "j@j.com", "11999", time.Now(), time.Now())
	repo := mocks.NewMockRequesterRepository(ctrl)
	repo.EXPECT().Update(gomock.Any(), c).Return(nil)

	uc := NewUpdateRequester(repo)
	if err := uc.Execute(context.Background(), c); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestUpdateRequester_Execute_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoErr := errors.New("db error")
	c := entities.ReconstituteRequester("", "João", "123", "j@j.com", "11999", time.Now(), time.Now())
	repo := mocks.NewMockRequesterRepository(ctrl)
	repo.EXPECT().Update(gomock.Any(), c).Return(repoErr)

	uc := NewUpdateRequester(repo)
	if err := uc.Execute(context.Background(), c); !errors.Is(err, repoErr) {
		t.Fatalf("expected %v, got %v", repoErr, err)
	}
}
