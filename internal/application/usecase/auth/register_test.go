package authuc

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestRegister_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	hasher := mocks.NewMockPasswordHasher(ctrl)

	userRepo.EXPECT().FindByEmail(gomock.Any(), "user@email.com").Return(nil, domainerrors.ErrNotFound)
	hasher.EXPECT().Hash("senha123").Return("hashed", nil)
	userRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	uc := NewRegister(userRepo, hasher)
	err := uc.Execute(context.Background(), "User", "user@email.com", "senha123", entities.RoleClient)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestRegister_Execute_EmailAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	existing := entities.NewUser("User", "user@email.com", "hash", entities.RoleClient)
	userRepo := mocks.NewMockUserRepository(ctrl)
	hasher := mocks.NewMockPasswordHasher(ctrl)

	userRepo.EXPECT().FindByEmail(gomock.Any(), "user@email.com").Return(existing, nil)

	uc := NewRegister(userRepo, hasher)
	err := uc.Execute(context.Background(), "User", "user@email.com", "senha123", entities.RoleClient)
	if !errors.Is(err, domainerrors.ErrAlreadyExists) {
		t.Fatalf("expected ErrAlreadyExists, got %v", err)
	}
}

func TestRegister_Execute_FindByEmailError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	hasher := mocks.NewMockPasswordHasher(ctrl)
	repoErr := errors.New("db error")

	userRepo.EXPECT().FindByEmail(gomock.Any(), "user@email.com").Return(nil, repoErr)

	uc := NewRegister(userRepo, hasher)
	err := uc.Execute(context.Background(), "User", "user@email.com", "senha123", entities.RoleClient)
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected db error, got %v", err)
	}
}

func TestRegister_Execute_HashError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	hasher := mocks.NewMockPasswordHasher(ctrl)
	hashErr := errors.New("hash error")

	userRepo.EXPECT().FindByEmail(gomock.Any(), "user@email.com").Return(nil, domainerrors.ErrNotFound)
	hasher.EXPECT().Hash("senha123").Return("", hashErr)

	uc := NewRegister(userRepo, hasher)
	err := uc.Execute(context.Background(), "User", "user@email.com", "senha123", entities.RoleClient)
	if !errors.Is(err, hashErr) {
		t.Fatalf("expected hash error, got %v", err)
	}
}
