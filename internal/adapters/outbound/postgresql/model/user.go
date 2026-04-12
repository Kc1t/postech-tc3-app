package pgmodel

import (
	"time"

	"github.com/fiap/postech-tc1/internal/domain/entities"
)

type User struct {
	ID             string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name           string     `gorm:"not null"`
	Email          string     `gorm:"uniqueIndex;not null"`
	PasswordHash   string     `gorm:"not null"`
	Role           string     `gorm:"not null;default:'client'"`
	FailedAttempts int        `gorm:"not null;default:0"`
	LockedUntil    *time.Time // nullable
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func FromUser(u *entities.User) *User {
	return &User{
		ID:             u.ID(),
		Name:           u.Name(),
		Email:          u.Email(),
		PasswordHash:   u.PasswordHash(),
		Role:           string(u.Role()),
		FailedAttempts: u.FailedAttempts(),
		LockedUntil:    u.LockedUntil(),
		CreatedAt:      u.CreatedAt(),
		UpdatedAt:      u.UpdatedAt(),
	}
}

func (m *User) ToDomain() *entities.User {
	return entities.ReconstituteUser(
		m.ID, m.Name, m.Email, m.PasswordHash,
		entities.Role(m.Role),
		m.FailedAttempts, m.LockedUntil,
		m.CreatedAt, m.UpdatedAt,
	)
}
