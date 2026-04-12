package postgresql

import (
	"context"
	"errors"

	pgmodel "github.com/fiap/postech-tc1/internal/adapters/outbound/postgresql/model"
	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/ports"
	"gorm.io/gorm"
)

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) ports.RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(ctx context.Context, rt *entities.RefreshToken) error {
	m := pgmodel.FromRefreshToken(rt)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	rt.SetID(m.ID)
	return nil
}

func (r *refreshTokenRepository) FindByTokenHash(ctx context.Context, hash string) (*entities.RefreshToken, error) {
	var m pgmodel.RefreshToken
	if err := r.db.WithContext(ctx).First(&m, "token_hash = ? AND revoked = false", hash).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrNotFound
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *refreshTokenRepository) Revoke(ctx context.Context, id string) error {
	// Condicao revoked = false garante single-use real em cenarios concorrentes:
	// apenas uma transacao consegue marcar o token como revogado.
	result := r.db.WithContext(ctx).
		Model(&pgmodel.RefreshToken{}).
		Where("id = ? AND revoked = false", id).
		Update("revoked", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainerrors.ErrNotFound
	}
	return nil
}

func (r *refreshTokenRepository) RevokeByUserID(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).
		Model(&pgmodel.RefreshToken{}).
		Where("user_id = ? AND revoked = false", userID).
		Update("revoked", true).Error
}

// RotateToken atomicamente revoga o token antigo e cria um novo dentro de uma transacao.
func (r *refreshTokenRepository) RotateToken(ctx context.Context, oldTokenHash string, newRT *entities.RefreshToken) (*entities.RefreshToken, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Buscar o token antigo
		var old pgmodel.RefreshToken
		if err := tx.First(&old, "token_hash = ? AND revoked = false", oldTokenHash).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domainerrors.ErrInvalidRefreshToken
			}
			return err
		}

		// 2. Validar
		domain := old.ToDomain()
		if !domain.IsValid() {
			return domainerrors.ErrInvalidRefreshToken
		}

		// 3. Revogar (single-use via WHERE revoked = false)
		result := tx.Model(&pgmodel.RefreshToken{}).
			Where("id = ? AND revoked = false", old.ID).
			Update("revoked", true)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return domainerrors.ErrInvalidRefreshToken
		}

		// 4. Criar novo refresh token
		m := pgmodel.FromRefreshToken(newRT)
		if err := tx.Create(m).Error; err != nil {
			return err
		}
		newRT.SetID(m.ID)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return newRT, nil
}
