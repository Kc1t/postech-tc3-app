	package authuc

import (
	"context"
	"errors"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/ports"
)

type Register struct {
	userRepo ports.UserRepository
	hasher   ports.PasswordHasher
}

func NewRegister(userRepo ports.UserRepository, hasher ports.PasswordHasher) *Register {
	return &Register{userRepo: userRepo, hasher: hasher}
}

func (uc *Register) Execute(ctx context.Context, name, email, rawPassword string, role entities.Role) error {
	existing, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, domainerrors.ErrNotFound) {
		return err
	}
	if existing != nil {
		return domainerrors.ErrAlreadyExists
	}

	hash, err := uc.hasher.Hash(rawPassword)
	if err != nil {
		return err
	}

	u := entities.NewUser(name, email, hash, role)
	return uc.userRepo.Create(ctx, u)
}
