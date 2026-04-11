package pgmodel

import (
	"time"

	"github.com/fiap/postech-tc1/internal/domain/user"
)

type RefreshToken struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    string    `gorm:"index;not null"`
	TokenHash string    `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	Revoked   bool      `gorm:"default:false"`
	CreatedAt time.Time
}

func FromRefreshToken(rt *user.RefreshToken) *RefreshToken {
	return &RefreshToken{
		ID:        rt.ID(),
		UserID:    rt.UserID(),
		TokenHash: rt.TokenHash(),
		ExpiresAt: rt.ExpiresAt(),
		Revoked:   rt.Revoked(),
		CreatedAt: rt.CreatedAt(),
	}
}

func (m *RefreshToken) ToDomain() *user.RefreshToken {
	return user.ReconstituteRefreshToken(m.ID, m.UserID, m.TokenHash, m.ExpiresAt, m.Revoked, m.CreatedAt)
}
