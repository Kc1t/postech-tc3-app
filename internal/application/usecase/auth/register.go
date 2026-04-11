package authuc

import (
	"context"
	"errors"

	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/domain/user"
	"github.com/fiap/postech-tc1/internal/ports"
	"golang.org/x/crypto/bcrypt"
)

type Register struct {
	userRepo   ports.UserRepository
	bcryptCost int
}

func NewRegister(userRepo ports.UserRepository, bcryptCost int) *Register {
	return &Register{userRepo: userRepo, bcryptCost: bcryptCost}
}

func (uc *Register) Execute(ctx context.Context, name, email, rawPassword string, role user.Role) error {
	// 1. Verificar se email ja existe
	existing, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, domainerrors.ErrNotFound) {
		return err
	}
	if existing != nil {
		return domainerrors.ErrAlreadyExists
	}

	// 2. Hash da senha (logica de aplicacao, nao de handler)
	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), uc.bcryptCost)
	if err != nil {
		return err
	}

	// 3. Criar e persistir
	u := user.New(name, email, string(hash), role)
	return uc.userRepo.Create(ctx, u)
}
