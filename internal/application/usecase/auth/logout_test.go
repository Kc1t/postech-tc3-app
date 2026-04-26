package authuc

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestLogout_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	refreshRepo := mocks.NewMockRefreshTokenRepository(ctrl)
	refreshRepo.EXPECT().RevokeByUserID(gomock.Any(), "user-1").Return(nil)

	uc := NewLogout(refreshRepo)
	if err := uc.Execute(context.Background(), "user-1"); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestLogout_Execute_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoErr := errors.New("db error")
	refreshRepo := mocks.NewMockRefreshTokenRepository(ctrl)
	refreshRepo.EXPECT().RevokeByUserID(gomock.Any(), "user-1").Return(repoErr)

	uc := NewLogout(refreshRepo)
	if err := uc.Execute(context.Background(), "user-1"); !errors.Is(err, repoErr) {
		t.Fatalf("expected db error, got %v", err)
	}
}
