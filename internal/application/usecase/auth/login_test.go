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

func newLoginUC(t *testing.T, ctrl *gomock.Controller) (*Login, *mocks.MockUserRepository, *mocks.MockRefreshTokenRepository, *mocks.MockTokenService, *mocks.MockPasswordHasher) {
	t.Helper()
	userRepo := mocks.NewMockUserRepository(ctrl)
	refreshRepo := mocks.NewMockRefreshTokenRepository(ctrl)
	tokenSvc := mocks.NewMockTokenService(ctrl)
	hasher := mocks.NewMockPasswordHasher(ctrl)

	hasher.EXPECT().Hash("timing-safe-dummy").Return("dummy-hash", nil)
	uc, err := NewLogin(userRepo, refreshRepo, tokenSvc, hasher, 5, 15)
	if err != nil {
		t.Fatalf("NewLogin: %v", err)
	}
	return uc, userRepo, refreshRepo, tokenSvc, hasher
}

func TestLogin_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc, userRepo, refreshRepo, tokenSvc, hasher := newLoginUC(t, ctrl)

	u := entities.NewUser("User", "user@email.com", "hashed", entities.RoleClient)
	userRepo.EXPECT().FindByEmail(gomock.Any(), "user@email.com").Return(u, nil)
	hasher.EXPECT().Compare("hashed", "senha123").Return(nil)
	userRepo.EXPECT().Update(gomock.Any(), u).Return(nil)
	tokenSvc.EXPECT().GenerateAccessToken(u).Return("access-token", nil)
	tokenSvc.EXPECT().GenerateRefreshToken().Return("raw-refresh", "hash-refresh", nil)
	tokenSvc.EXPECT().RefreshTokenExpiration().Return(7 * 24 * time.Hour)
	refreshRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	access, refresh, err := uc.Execute(context.Background(), "user@email.com", "senha123")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if access != "access-token" || refresh != "raw-refresh" {
		t.Fatalf("unexpected tokens: %s %s", access, refresh)
	}
}

func TestLogin_Execute_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc, userRepo, _, _, hasher := newLoginUC(t, ctrl)

	userRepo.EXPECT().FindByEmail(gomock.Any(), "notfound@email.com").Return(nil, domainerrors.ErrNotFound)
	hasher.EXPECT().Compare("dummy-hash", "senha123").Return(errors.New("mismatch"))

	_, _, err := uc.Execute(context.Background(), "notfound@email.com", "senha123")
	if !errors.Is(err, domainerrors.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLogin_Execute_WrongPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc, userRepo, _, _, hasher := newLoginUC(t, ctrl)

	u := entities.NewUser("User", "user@email.com", "hashed", entities.RoleClient)
	userRepo.EXPECT().FindByEmail(gomock.Any(), "user@email.com").Return(u, nil)
	hasher.EXPECT().Compare("hashed", "wrong").Return(errors.New("mismatch"))
	userRepo.EXPECT().Update(gomock.Any(), u).Return(nil)

	_, _, err := uc.Execute(context.Background(), "user@email.com", "wrong")
	if !errors.Is(err, domainerrors.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLogin_Execute_LockedAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc, userRepo, _, _, hasher := newLoginUC(t, ctrl)

	locked := time.Now().Add(15 * time.Minute)
	u := entities.ReconstituteUser("id-1", "User", "user@email.com", "hashed", entities.RoleClient, 5, &locked, time.Now(), time.Now())
	userRepo.EXPECT().FindByEmail(gomock.Any(), "user@email.com").Return(u, nil)
	hasher.EXPECT().Compare("hashed", "senha123").Return(nil)

	_, _, err := uc.Execute(context.Background(), "user@email.com", "senha123")
	if !errors.Is(err, domainerrors.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLogin_Execute_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc, userRepo, _, _, _ := newLoginUC(t, ctrl)

	repoErr := errors.New("db error")
	userRepo.EXPECT().FindByEmail(gomock.Any(), "user@email.com").Return(nil, repoErr)

	_, _, err := uc.Execute(context.Background(), "user@email.com", "senha123")
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected db error, got %v", err)
	}
}
