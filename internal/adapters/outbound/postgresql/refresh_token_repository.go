package postgresql

import (
	"context"
	"errors"

	pgmodel "github.com/fiap/postech-tc1/internal/adapters/outbound/postgresql/model"
	"github.com/fiap/postech-tc1/internal/domain/user"
	"github.com/fiap/postech-tc1/internal/ports"
	"gorm.io/gorm"
)

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) ports.RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(ctx context.Context, rt *user.RefreshToken) error {
	db := GetDB(ctx, r.db)
	m := pgmodel.FromRefreshToken(rt)
	if err := db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	rt.SetID(m.ID)
	return nil
}

func (r *refreshTokenRepository) FindByTokenHash(ctx context.Context, hash string) (*user.RefreshToken, error) {
	db := GetDB(ctx, r.db)
	var m pgmodel.RefreshToken
	if err := db.WithContext(ctx).First(&m, "token_hash = ? AND revoked = false", hash).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *refreshTokenRepository) Revoke(ctx context.Context, id string) error {
	db := GetDB(ctx, r.db)
	return db.WithContext(ctx).
		Model(&pgmodel.RefreshToken{}).
		Where("id = ?", id).
		Update("revoked", true).Error
}

func (r *refreshTokenRepository) RevokeByUserID(ctx context.Context, userID string) error {
	db := GetDB(ctx, r.db)
	return db.WithContext(ctx).
		Model(&pgmodel.RefreshToken{}).
		Where("user_id = ? AND revoked = false", userID).
		Update("revoked", true).Error
}
