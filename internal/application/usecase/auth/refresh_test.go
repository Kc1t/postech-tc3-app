package authuc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestRefresh_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	refreshRepo := mocks.NewMockRefreshTokenRepository(ctrl)
	tokenSvc := mocks.NewMockTokenService(ctrl)

	rt := entities.NewRefreshToken("user-1", "hashed-token", time.Now().Add(time.Hour))
	u := entities.NewUser("User", "user@email.com", "hash", entities.RoleClient)

	tokenSvc.EXPECT().HashToken("raw-token").Return("hashed-token")
	refreshRepo.EXPECT().FindByTokenHash(gomock.Any(), "hashed-token").Return(rt, nil)
	userRepo.EXPECT().FindByID(gomock.Any(), "user-1").Return(u, nil)
	tokenSvc.EXPECT().GenerateAccessToken(u).Return("new-access", nil)
	tokenSvc.EXPECT().GenerateRefreshToken().Return("new-raw", "new-hash", nil)
	tokenSvc.EXPECT().RefreshTokenExpiration().Return(7 * 24 * time.Hour)
	refreshRepo.EXPECT().RotateToken(gomock.Any(), rt.ID(), gomock.Any()).Return(nil)

	uc := NewRefresh(userRepo, refreshRepo, tokenSvc)
	access, refresh, err := uc.Execute(context.Background(), "raw-token")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if access != "new-access" || refresh != "new-raw" {
		t.Fatalf("unexpected tokens: %s %s", access, refresh)
	}
}

func TestRefresh_Execute_TokenNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	refreshRepo := mocks.NewMockRefreshTokenRepository(ctrl)
	tokenSvc := mocks.NewMockTokenService(ctrl)

	tokenSvc.EXPECT().HashToken("raw-token").Return("hashed-token")
	refreshRepo.EXPECT().FindByTokenHash(gomock.Any(), "hashed-token").Return(nil, domainerrors.ErrNotFound)

	uc := NewRefresh(userRepo, refreshRepo, tokenSvc)
	_, _, err := uc.Execute(context.Background(), "raw-token")
	if !errors.Is(err, domainerrors.ErrInvalidRefreshToken) {
		t.Fatalf("expected ErrInvalidRefreshToken, got %v", err)
	}
}

func TestRefresh_Execute_ExpiredToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	refreshRepo := mocks.NewMockRefreshTokenRepository(ctrl)
	tokenSvc := mocks.NewMockTokenService(ctrl)

	expired := entities.NewRefreshToken("user-1", "hashed-token", time.Now().Add(-time.Hour))
	tokenSvc.EXPECT().HashToken("raw-token").Return("hashed-token")
	refreshRepo.EXPECT().FindByTokenHash(gomock.Any(), "hashed-token").Return(expired, nil)

	uc := NewRefresh(userRepo, refreshRepo, tokenSvc)
	_, _, err := uc.Execute(context.Background(), "raw-token")
	if !errors.Is(err, domainerrors.ErrInvalidRefreshToken) {
		t.Fatalf("expected ErrInvalidRefreshToken, got %v", err)
	}
}

func TestRefresh_Execute_RotateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	refreshRepo := mocks.NewMockRefreshTokenRepository(ctrl)
	tokenSvc := mocks.NewMockTokenService(ctrl)

	rt := entities.NewRefreshToken("user-1", "hashed-token", time.Now().Add(time.Hour))
	u := entities.NewUser("User", "user@email.com", "hash", entities.RoleClient)

	tokenSvc.EXPECT().HashToken("raw-token").Return("hashed-token")
	refreshRepo.EXPECT().FindByTokenHash(gomock.Any(), "hashed-token").Return(rt, nil)
	userRepo.EXPECT().FindByID(gomock.Any(), "user-1").Return(u, nil)
	tokenSvc.EXPECT().GenerateAccessToken(u).Return("new-access", nil)
	tokenSvc.EXPECT().GenerateRefreshToken().Return("new-raw", "new-hash", nil)
	tokenSvc.EXPECT().RefreshTokenExpiration().Return(7 * 24 * time.Hour)
	refreshRepo.EXPECT().RotateToken(gomock.Any(), rt.ID(), gomock.Any()).Return(domainerrors.ErrNotFound)

	uc := NewRefresh(userRepo, refreshRepo, tokenSvc)
	_, _, err := uc.Execute(context.Background(), "raw-token")
	if !errors.Is(err, domainerrors.ErrInvalidRefreshToken) {
		t.Fatalf("expected ErrInvalidRefreshToken, got %v", err)
	}
}
