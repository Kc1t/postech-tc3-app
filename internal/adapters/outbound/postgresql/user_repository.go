package postgresql

import (
	"context"
	"errors"

	pgmodel "github.com/fiap/postech-tc1/internal/adapters/outbound/postgresql/model"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/domain/user"
	"github.com/fiap/postech-tc1/internal/ports"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) ports.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u *user.User) error {
	db := GetDB(ctx, r.db)
	m := pgmodel.FromUser(u)
	if err := db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	u.SetID(m.ID)
	return nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	db := GetDB(ctx, r.db)
	var m pgmodel.User
	if err := db.WithContext(ctx).First(&m, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrNotFound
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	db := GetDB(ctx, r.db)
	var m pgmodel.User
	if err := db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrNotFound
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *userRepository) Update(ctx context.Context, u *user.User) error {
	db := GetDB(ctx, r.db)
	m := pgmodel.FromUser(u)
	// Save atualiza todos os campos — usado pra persistir mudancas de lockout
	// (failed_attempts, locked_until) e dados do perfil.
	return db.WithContext(ctx).Save(m).Error
}
